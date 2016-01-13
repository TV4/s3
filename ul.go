package s3

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Upload(bucket, id, secret, region, key string, data []byte) {
	ul := setupBurtUploader(bucket, id, secret, region)
	ul.upload(data, key)
}

type s3Uploader struct {
	bucket   string
	s3Client *s3.S3
	uloader  *s3manager.Uploader
}

func setupBurtUploader(bucket, id, secret, region string) s3Uploader {
	s3Client := s3.New(session.New(), awsConfig(id, secret, region))

	uloader := &s3manager.Uploader{
		PartSize:       1024 * 1024 * 20,
		MaxUploadParts: 100,
		Concurrency:    15,
		S3:             s3Client,
	}

	return s3Uploader{
		bucket:   bucket,
		s3Client: s3Client,
		uloader:  uloader,
	}
}

func (ul *s3Uploader) upload(data []byte, key string) {
	params := &s3manager.UploadInput{
		Bucket: &ul.bucket,
		Key:    &key,
		Body:   bytes.NewBuffer(data),
	}
	fmt.Println(ul.bucket, key)
	res, err := ul.uloader.Upload(params)
	if err != nil {
		fmt.Println("error uploading:", err)
		return
	}
	fmt.Println("successfully uploaded:", res.Location)
}
