package s3lib

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	// Test configuration using environment variables or defaults
	testConfig = Config{
		Region:    getEnvOrDefault("AWS_REGION", "us-west-2"),
		AccessKey: getEnvOrDefault("AWS_ACCESS_KEY_ID", "test-key"),
		SecretKey: getEnvOrDefault("AWS_SECRET_ACCESS_KEY", "test-secret"),
		Duration:  5 * time.Minute,
		Debug:     true,
	}

	// Test constants
	testBucket      = getEnvOrDefault("TEST_BUCKET", "test-bucket")
	testFileName    = "test-file.txt"
	testFileContent = []byte("Hello, World!")
)

// Helper function to get environment variables with defaults
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper function to setup test client
func setupTestClient(t *testing.T) *S3Client {
	client, err := NewS3Client(testConfig)
	require.NoError(t, err)
	require.NotNil(t, client)
	return client
}

// TestNewS3Client tests the creation of a new S3 client
func TestNewS3Client(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "Valid configuration",
			config:  testConfig,
			wantErr: false,
		},
		{
			name: "Empty region",
			config: Config{
				Region:    "",
				AccessKey: "test-key",
				SecretKey: "test-secret",
			},
			wantErr: true,
		},
		{
			name: "Empty access key",
			config: Config{
				Region:    "us-west-2",
				AccessKey: "",
				SecretKey: "test-secret",
			},
			wantErr: true,
		},
		{
			name: "Empty secret key",
			config: Config{
				Region:    "us-west-2",
				AccessKey: "test-key",
				SecretKey: "",
			},
			wantErr: true,
		},
		{
			name: "With endpoint",
			config: Config{
				Region:    "us-west-2",
				AccessKey: "test-key",
				SecretKey: "test-secret",
				Endpoint:  "http://localhost:4566", // LocalStack endpoint
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewS3Client(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

// TestS3Client_ListFiles tests the ListFiles function
func TestS3Client_ListFiles(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		name    string
		bucket  string
		prefix  string
		wantErr bool
	}{
		{
			name:    "Valid bucket",
			bucket:  testBucket,
			prefix:  "",
			wantErr: false,
		},
		{
			name:    "With prefix",
			bucket:  testBucket,
			prefix:  "test/",
			wantErr: false,
		},
		{
			name:    "Empty bucket",
			bucket:  "",
			prefix:  "",
			wantErr: true,
		},
		{
			name:    "Invalid bucket",
			bucket:  "nonexistent-bucket",
			prefix:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			files, err := client.ListFiles(ctx, tt.bucket, tt.prefix)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, files)
			}
		})
	}
}

// TestS3Client_UploadFile tests the UploadFile function
func TestS3Client_UploadFile(t *testing.T) {
	client := setupTestClient(t)

	tests := []struct {
		name     string
		bucket   string
		filename string
		data     []byte
		opts     *UploadOptions
		wantErr  bool
	}{
		{
			name:     "Valid upload",
			bucket:   testBucket,
			filename: testFileName,
			data:     testFileContent,
			opts:     nil,
			wantErr:  false,
		},
		{
			name:     "With options",
			bucket:   testBucket,
			filename: "test-with-opts.txt",
			data:     testFileContent,
			opts: &UploadOptions{
				ContentType: "text/plain",
				Metadata: map[string]string{
					"test": "value",
				},
			},
			wantErr: false,
		},
		{
			name:     "Empty bucket",
			bucket:   "",
			filename: testFileName,
			data:     testFileContent,
			opts:     nil,
			wantErr:  true,
		},
		{
			name:     "Empty filename",
			bucket:   testBucket,
			filename: "",
			data:     testFileContent,
			opts:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			location, err := client.UploadFile(ctx, tt.bucket, tt.filename, tt.data, tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, location)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, location)
			}
		})
	}
}

// TestS3Client_DownloadFile tests the DownloadFile function
func TestS3Client_DownloadFile(t *testing.T) {
	client := setupTestClient(t)

	// Upload a test file first
	ctx := context.Background()
	_, err := client.UploadFile(ctx, testBucket, testFileName, testFileContent, nil)
	require.NoError(t, err)

	tests := []struct {
		name    string
		bucket  string
		key     string
		wantErr bool
	}{
		{
			name:    "Valid download",
			bucket:  testBucket,
			key:     testFileName,
			wantErr: false,
		},
		{
			name:    "Non-existent file",
			bucket:  testBucket,
			key:     "nonexistent.txt",
			wantErr: true,
		},
		{
			name:    "Empty bucket",
			bucket:  "",
			key:     testFileName,
			wantErr: true,
		},
		{
			name:    "Empty key",
			bucket:  testBucket,
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := client.DownloadFile(ctx, tt.bucket, tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testFileContent, data)
			}
		})
	}
}

// TestS3Client_DeleteFile tests the DeleteFile function
func TestS3Client_DeleteFile(t *testing.T) {
	client := setupTestClient(t)

	// Upload a test file first
	ctx := context.Background()
	_, err := client.UploadFile(ctx, testBucket, "to-delete.txt", testFileContent, nil)
	require.NoError(t, err)

	tests := []struct {
		name    string
		bucket  string
		key     string
		wantErr bool
	}{
		{
			name:    "Valid delete",
			bucket:  testBucket,
			key:     "to-delete.txt",
			wantErr: false,
		},
		{
			name:    "Empty bucket",
			bucket:  "",
			key:     "test.txt",
			wantErr: true,
		},
		{
			name:    "Empty key",
			bucket:  testBucket,
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteFile(ctx, tt.bucket, tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestS3Client_GetFileInfo tests the GetFileInfo function
func TestS3Client_GetFileInfo(t *testing.T) {
	client := setupTestClient(t)

	// Upload a test file first
	ctx := context.Background()
	_, err := client.UploadFile(ctx, testBucket, testFileName, testFileContent, nil)
	require.NoError(t, err)

	tests := []struct {
		name    string
		bucket  string
		key     string
		wantErr bool
	}{
		{
			name:    "Valid file info",
			bucket:  testBucket,
			key:     testFileName,
			wantErr: false,
		},
		{
			name:    "Non-existent file",
			bucket:  testBucket,
			key:     "nonexistent.txt",
			wantErr: true,
		},
		{
			name:    "Empty bucket",
			bucket:  "",
			key:     testFileName,
			wantErr: true,
		},
		{
			name:    "Empty key",
			bucket:  testBucket,
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := client.GetFileInfo(ctx, tt.bucket, tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, info)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, tt.key, info.Key)
				assert.Equal(t, int64(len(testFileContent)), info.Size)
			}
		})
	}
}

// TestIntegration performs an end-to-end test
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	ctx := context.Background()

	// Test full workflow
	t.Run("Full workflow", func(t *testing.T) {
		// 1. Upload file
		filename := "integration-test.txt"
		content := []byte("Integration test content")
		location, err := client.UploadFile(ctx, testBucket, filename, content, &UploadOptions{
			ContentType: "text/plain",
			Metadata: map[string]string{
				"test": "integration",
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, location)

		// 2. Get file info
		info, err := client.GetFileInfo(ctx, testBucket, filename)
		require.NoError(t, err)
		require.Equal(t, filename, info.Key)
		require.Equal(t, int64(len(content)), info.Size)

		// 3. Download file
		downloaded, err := client.DownloadFile(ctx, testBucket, filename)
		require.NoError(t, err)
		require.Equal(t, content, downloaded)

		// 4. List files
		files, err := client.ListFiles(ctx, testBucket, "")
		require.NoError(t, err)
		found := false
		for _, file := range files {
			if file.Key == filename {
				found = true
				break
			}
		}
		require.True(t, found)

		// 5. Delete file
		err = client.DeleteFile(ctx, testBucket, filename)
		require.NoError(t, err)

		// 6. Verify deletion
		_, err = client.GetFileInfo(ctx, testBucket, filename)
		require.Error(t, err)
	})
}

// Examples
func ExampleS3Client_UploadFile() {
	cfg := Config{
		Region:    "us-west-2",
		AccessKey: "your-access-key",
		SecretKey: "your-secret-key",
		Duration:  5 * time.Minute,
	}

	client, err := NewS3Client(cfg)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	ctx := context.Background()
	data := []byte("Hello, World!")
	location, err := client.UploadFile(ctx, "example-bucket", "example.txt", data, nil)
	if err != nil {
		fmt.Printf("Failed to upload: %v\n", err)
		return
	}

	fmt.Printf("File uploaded to: %s\n", location)
}

func ExampleS3Client_DownloadFile() {
	cfg := Config{
		Region:    "us-west-2",
		AccessKey: "your-access-key",
		SecretKey: "your-secret-key",
		Duration:  5 * time.Minute,
	}

	client, err := NewS3Client(cfg)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	ctx := context.Background()
	data, err := client.DownloadFile(ctx, "example-bucket", "example.txt")
	if err != nil {
		fmt.Printf("Failed to download: %v\n", err)
		return
	}

	fmt.Printf("Downloaded content: %s\n", string(data))
}
