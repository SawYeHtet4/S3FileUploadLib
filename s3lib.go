package s3lib

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Client represents the S3 client configuration and operations
type S3Client struct {
	s3Client  *s3.S3
	session   *session.Session
	uploader  *s3manager.Uploader
	config    Config
	debugMode bool
}

// FileInfo represents S3 object metadata
type FileInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	StorageClass string    `json:"storage_class"`
}

// UploadOptions represents optional parameters for upload operations
type UploadOptions struct {
	ContentType        string
	ContentDisposition string
	CacheControl       string
	Metadata           map[string]string
	StorageClass       string
	ACL                string
}

// NewS3Client creates a new S3 client instance
func NewS3Client(cfg Config) (*S3Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	awsCfg := &aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
	}

	if cfg.Endpoint != "" {
		awsCfg.Endpoint = aws.String(cfg.Endpoint)
		awsCfg.S3ForcePathStyle = aws.Bool(true)
	}

	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	s3Client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &S3Client{
		s3Client:  s3Client,
		session:   sess,
		uploader:  uploader,
		config:    cfg,
		debugMode: cfg.Debug,
	}, nil
}

// ListFiles lists all files in the specified bucket with optional prefix
func (c *S3Client) ListFiles(ctx context.Context, bucket, prefix string) ([]FileInfo, error) {
	if bucket == "" {
		return nil, ErrInvalidBucket
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	var files []FileInfo
	err := c.s3Client.ListObjectsV2PagesWithContext(ctx, input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				files = append(files, FileInfo{
					Key:          aws.StringValue(obj.Key),
					Size:         aws.Int64Value(obj.Size),
					LastModified: aws.TimeValue(obj.LastModified),
					ETag:         aws.StringValue(obj.ETag),
					StorageClass: aws.StringValue(obj.StorageClass),
				})
			}
			return true
		})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, ErrInvalidBucket
			default:
				return nil, fmt.Errorf("AWS error: %w", aerr)
			}
		}
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	return files, nil
}

// UploadFile uploads a file to the specified bucket with options
func (c *S3Client) UploadFile(ctx context.Context, bucket, filename string, data []byte, opts *UploadOptions) (string, error) {
	if bucket == "" {
		return "", ErrInvalidBucket
	}
	if filename == "" {
		return "", ErrInvalidKey
	}

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	}

	if opts != nil {
		if opts.ContentType != "" {
			input.ContentType = aws.String(opts.ContentType)
		}
		if opts.ContentDisposition != "" {
			input.ContentDisposition = aws.String(opts.ContentDisposition)
		}
		if opts.CacheControl != "" {
			input.CacheControl = aws.String(opts.CacheControl)
		}
		if opts.Metadata != nil {
			input.Metadata = aws.StringMap(opts.Metadata)
		}
		if opts.StorageClass != "" {
			input.StorageClass = aws.String(opts.StorageClass)
		}
		if opts.ACL != "" {
			input.ACL = aws.String(opts.ACL)
		}
	}

	result, err := c.uploader.UploadWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return "", ErrInvalidBucket
			default:
				return "", fmt.Errorf("AWS error: %w", aerr)
			}
		}
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return result.Location, nil
}

// DownloadFile downloads a file from the specified bucket
func (c *S3Client) DownloadFile(ctx context.Context, bucket, key string) ([]byte, error) {
	if bucket == "" {
		return nil, ErrInvalidBucket
	}
	if key == "" {
		return nil, ErrInvalidKey
	}

	// First check if the object exists
	_, err := c.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return nil, ErrFileNotFound
			case s3.ErrCodeNoSuchBucket:
				return nil, ErrInvalidBucket
			default:
				return nil, fmt.Errorf("AWS error: %w", aerr)
			}
		}
		return nil, fmt.Errorf("failed to get object info: %w", err)
	}

	// Download the object
	buf := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloader(c.session)

	_, err = downloader.DownloadWithContext(ctx, buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return buf.Bytes(), nil
}

// DeleteFile deletes a file from the specified bucket
func (c *S3Client) DeleteFile(ctx context.Context, bucket, key string) error {
	if bucket == "" {
		return ErrInvalidBucket
	}
	if key == "" {
		return ErrInvalidKey
	}

	_, err := c.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return ErrInvalidBucket
			default:
				return fmt.Errorf("AWS error: %w", aerr)
			}
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFileInfo gets metadata for a specific file
func (c *S3Client) GetFileInfo(ctx context.Context, bucket, key string) (*FileInfo, error) {
	if bucket == "" {
		return nil, ErrInvalidBucket
	}
	if key == "" {
		return nil, ErrInvalidKey
	}

	result, err := c.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return nil, ErrFileNotFound
			case s3.ErrCodeNoSuchBucket:
				return nil, ErrInvalidBucket
			default:
				return nil, fmt.Errorf("AWS error: %w", aerr)
			}
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &FileInfo{
		Key:          key,
		Size:         aws.Int64Value(result.ContentLength),
		LastModified: aws.TimeValue(result.LastModified),
		ETag:         aws.StringValue(result.ETag),
		StorageClass: aws.StringValue(result.StorageClass),
	}, nil
}

// Close closes the S3 client and cleans up resources
func (c *S3Client) Close() error {
	if c == nil {
		return nil
	}

	// Cancel any ongoing operations
	if c.session != nil {
		// Clean up session-related resources
		c.session = nil
	}

	// Clean up the S3 client
	if c.s3Client != nil {
		c.s3Client = nil
	}

	// Clean up the uploader
	if c.uploader != nil {
		c.uploader = nil
	}

	// Log cleanup if debug mode is enabled
	if c.debugMode {
		fmt.Println("S3 client resources cleaned up")
	}

	return nil
}
