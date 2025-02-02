# S3Lib

A simple S3 operations library for Go applications.

## Installation

```bash
go get github.com/SawYeHtet4/S3FileUploadLib.git


Features
File upload to S3
File download from S3
List files in bucket
AWS credentials configuration
Context support for timeouts and cancellation
Usage

package main

import (
    "context"
    "log"
    "time"
    "github.com/yourusername/s3lib"
)

func main() {
    cfg := s3lib.Config{
        Region:    "us-west-2",
        AccessKey: "your-access-key",
        SecretKey: "your-secret-key",
        Duration:  5 * time.Minute,
    }

    client, err := s3lib.NewS3Client(cfg)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Upload file
    ctx := context.Background()
    location, err := client.UploadFile(ctx, "my-bucket", "test.txt", []byte("Hello, World!"))
    if err != nil {
        log.Printf("Failed to upload: %v", err)
        return
    }
    log.Printf("File uploaded to: %s", location)
}
