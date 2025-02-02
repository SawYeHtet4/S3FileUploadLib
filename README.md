# S3Lib

A simple S3 operations library for Go applications.

_Features_
   1. File upload to S3
   2. File download from S3
   3. List files in bucket
   4. AWS credentials configuration
   5. Context support for timeouts and cancellation

## Installation

```bash
go get github.com/SawYeHtet4/S3FileUploadLib.git
```

## Usage

```bash
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
```
