package storage

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	bucket     string
	uploader   *s3manager.Uploader
	downloader *s3.S3
}

// NewS3Storage creates a new S3 storage client.
// Note on "bucket per user": Creating an S3 bucket per user is generally not recommended
// due to AWS account limits on the number of buckets and the slow speed of bucket creation.
// A better, more scalable approach is to use a single bucket with user-specific prefixes (folders),
// e.g., "user_123/post_abc.md". This implementation uses a single bucket name provided
// via configuration and assumes paths will contain user-specific identifiers.
func NewS3Storage(bucket, region string) (*S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		bucket:     bucket,
		uploader:   s3manager.NewUploader(sess),
		downloader: s3.New(sess),
	}, nil
}

func (s *S3Storage) Save(path string, data []byte) error {
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *S3Storage) Read(path string) ([]byte, error) {
	result, err := s.downloader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	return ioutil.ReadAll(result.Body)
}