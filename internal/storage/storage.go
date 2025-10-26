package storage

import (
	"fmt"
	"go-blog/internal/config"
)

// FileStorage defines the interface for file storage operations.
type FileStorage interface {
	Save(path string, data []byte) error
	Read(path string) ([]byte, error)
}

// New creates a new FileStorage instance based on the configuration.
func New(cfg *config.Config) (FileStorage, error) {
	switch cfg.StorageType {
	case "local":
		storage, err := NewLocalStorage("./files")
		if err != nil {
			return nil, err
		}
		return storage, nil
	case "s3":
		return NewS3Storage(cfg.S3Bucket, cfg.S3Region)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.StorageType)
	}
}