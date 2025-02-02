package s3lib

import "time"

// Config holds the configuration for S3Client
type Config struct {
    Region    string
    AccessKey string
    SecretKey string
    Duration  time.Duration
    Endpoint  string        // Optional: for S3-compatible services
    UseSSL    bool         // Optional: use HTTPS
    Debug     bool         // Optional: enable debug logging
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
    if c.Region == "" {
        return ErrInvalidConfig
    }
    if c.AccessKey == "" {
        return ErrInvalidConfig
    }
    if c.SecretKey == "" {
        return ErrInvalidConfig
    }
    return nil
}