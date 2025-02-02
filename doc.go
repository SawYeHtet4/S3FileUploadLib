/*
Package s3lib provides a simple interface for AWS S3 operations.

Basic usage:

    cfg := s3lib.Config{
        Region:    "us-west-2",
        AccessKey: "your-access-key",
        SecretKey: "your-secret-key",
        Duration:  5 * time.Minute,
    }

    client, err := s3lib.NewS3Client(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Upload a file
    data := []byte("Hello, World!")
    location, err := client.UploadFile(context.Background(), "my-bucket", "hello.txt", data, nil)

    // Download a file
    data, err := client.DownloadFile(context.Background(), "my-bucket", "hello.txt")

    // List files
    files, err := client.ListFiles(context.Background(), "my-bucket", "")

    // Delete a file
    err := client.DeleteFile(context.Background(), "my-bucket", "hello.txt")
*/
package s3lib