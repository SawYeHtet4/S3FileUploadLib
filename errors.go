package s3lib

import "errors"

var (
    // ErrInvalidConfig is returned when the configuration is invalid
    ErrInvalidConfig = errors.New("invalid configuration")
    
    // ErrInvalidBucket is returned when the bucket name is invalid
    ErrInvalidBucket = errors.New("invalid bucket name")
    
    // ErrInvalidKey is returned when the key is invalid
    ErrInvalidKey = errors.New("invalid key")
    
    // ErrFileNotFound is returned when the requested file is not found
    ErrFileNotFound = errors.New("file not found")
)