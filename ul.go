package s3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Upload starts the uploading of a binary object to the location given by
// path in the bucket specified by the configuration. It returns the object
// location in the bucket and an error value.
func Upload(c BucketConf, path string, data []byte) (string, error) {
	ul := NewUploader(c)
	return ul.Upload(data, path)
}

// Uploader is the interface that wraps the Upload method.
type Uploader interface {
	Upload(data []byte, path string) (string, error)
}

type s3Uploader struct {
	bucket   string
	s3Client *s3.S3
	uloader  *s3manager.Uploader
}

// NewUploader creates and initializes a new Uploader based on the values
// contained in BucketConf.
func NewUploader(c BucketConf) Uploader {
	s3Client := s3.New(session.New(), awsConfig(c.ID, c.Secret, c.Region))

	uloader := &s3manager.Uploader{
		PartSize:       1024 * 1024 * 20,
		MaxUploadParts: 100,
		Concurrency:    15,
		S3:             s3Client,
	}

	return s3Uploader{
		bucket:   c.Bucket,
		s3Client: s3Client,
		uloader:  uloader,
	}
}

func (u s3Uploader) Upload(data []byte, key string) (string, error) {
	params := &s3manager.UploadInput{
		Bucket: &u.bucket,
		Key:    &key,
		Body:   bytes.NewBuffer(data),
	}
	res, err := u.uloader.Upload(params)
	if err != nil {
		return "", err
	}

	return res.Location, nil
}
