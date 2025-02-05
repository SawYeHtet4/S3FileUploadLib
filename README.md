# S3Lib

A comprehensive S3 operations library for Go applications that provides simple interfaces for common S3 operations including direct uploads, pre-signed URLs, and pre-signed POST uploads.

## Features
1. Direct file operations:
   - Upload files to S3
   - Download files from S3
   - List files in bucket
   - Delete files from S3
   - Get file metadata
2. Pre-signed operations:
   - Generate pre-signed URLs for upload/download
   - Generate pre-signed POST data for browser uploads
3. Advanced features:
   - Custom metadata support
   - Content type handling
   - Storage class selection
   - Cache control settings
4. Security and configuration:
   - AWS credentials configuration
   - Context support for timeouts and cancellation
   - Debug mode for troubleshooting

## Installation

```bash
go get github.com/SawYeHtet4/S3FileUploadLib
```

## Usage

# Basic Setup

```bash
package main

import (
    "context"
    "log"
    "time"
    "github.com/SawYeHtet4/S3FileUploadLib"
)

func main() {
    cfg := s3lib.Config{
        Region:    "us-west-2",
        AccessKey: "your-access-key",
        SecretKey: "your-secret-key",
        Duration:  5 * time.Minute,
        Debug:     true, // Enable debug mode
    }

    client, err := s3lib.NewS3Client(cfg)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()
}
```

# Direct File Operations

```bash
// Upload file with options
uploadOpts := &s3lib.UploadOptions{
    ContentType:        "application/json",
    ContentDisposition: "attachment",
    CacheControl:       "max-age=3600",
    Metadata: map[string]string{
        "user-id": "123",
    },
    StorageClass: "STANDARD",
    ACL:          "private",
}

location, err := client.UploadFile(ctx, "my-bucket", "test.json", data, uploadOpts)

// Download file
data, err := client.DownloadFile(ctx, "my-bucket", "test.json")

// List files
files, err := client.ListFiles(ctx, "my-bucket", "prefix/")

// Delete file
err := client.DeleteFile(ctx, "my-bucket", "test.json")

// Get file info
info, err := client.GetFileInfo(ctx, "my-bucket", "test.json")
```


# Pre-signed URL Operations

```bash
// Generate pre-signed URL for upload
urlResp, err := client.GeneratePresignedURL(
    ctx,
    "my-bucket",
    "uploads/file.jpg",
    15*time.Minute,
    "upload",
)

// Generate pre-signed URL for download
urlResp, err := client.GeneratePresignedURL(
    ctx,
    "my-bucket",
    "downloads/file.jpg",
    15*time.Minute,
    "download",
)
```

# Pre-signed POST Operations

```bash 
// Generate pre-signed POST data
postResp, err := client.GeneratePresignedPost(
    ctx,
    "my-bucket",
    "uploads/file.jpg",
    15*time.Minute,
    10*1024*1024, // 10MB max size
)
```

# Client-side Upload Examples

```bash
// Using pre-signed URL
async function uploadWithURL(file) {
    const response = await fetch('/api/get-upload-url');
    const { url } = await response.json();
    
    await fetch(url, {
        method: 'PUT',
        body: file,
        headers: {
            'Content-Type': file.type
        }
    });
}

// Using pre-signed POST
async function uploadWithPost(file) {
    const response = await fetch('/api/get-post-data');
    const { url, fields } = await response.json();
    
    const formData = new FormData();
    Object.entries(fields).forEach(([key, value]) => {
        formData.append(key, value);
    });
    formData.append('file', file);
    
    await fetch(url, {
        method: 'POST',
        body: formData
    });
}
```

## Error Handling
# The library provides specific error types for common scenarios:

```bash
switch err {
    case s3lib.ErrInvalidConfig:
        // Hanle invalid configuration error
    case s3lib.ErrInvalidBucket:
        // Handle invalid bucket error
    case s3lib.ErrInvalidKey:
        // Handle invalid key error
    case s3lib.ErrFileNotFound:
        // Handle file not found error
    default:
        // Handle other errors
}
```

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.
