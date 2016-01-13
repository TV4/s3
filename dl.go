package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Downloader struct {
	bucket   string
	s3Client *s3.S3
	dloader  *s3manager.Downloader
}

type Config struct {
	Bucket, ID, Secret, Region string
}

func NewDownloader(c Config) s3Downloader {
	return NewS3Downloader(c.Bucket, c.ID, c.Secret, c.Region)
}

func NewS3Downloader(bucket, id, secret, region string) s3Downloader {
	s3Client := s3.New(session.New(), awsConfig(id, secret, region))
	dloader := &s3manager.Downloader{
		PartSize:    1024 * 1024 * 5,
		Concurrency: 5,
		S3:          s3Client,
	}

	return s3Downloader{
		bucket:   bucket,
		s3Client: s3Client,
		dloader:  dloader,
	}
}

func (dldr *s3Downloader) Download(path string, handler ChunkHandler) (<-chan int, <-chan error) {
	return dldr.DownloadChunks(path, handler, 0)
}

func (dldr *s3Downloader) DownloadChunks(path string, handler ChunkHandler, chunks int) (<-chan int, <-chan error) {
	objPath := &s3.ListObjectsInput{Bucket: &dldr.bucket, Prefix: &path}
	cntc, errc := make(chan int), make(chan error)

	go func() {
		// Iterate over objects located in objPath
		err := dldr.s3Client.ListObjectsPages(objPath, func(p *s3.ListObjectsOutput, lastPage bool) bool {
			if chunks <= 0 || chunks > len(p.Contents) {
				chunks = len(p.Contents)
			}
			cntc <- chunks
			for i := 0; i < chunks; i++ {
				obj, err := dldr.downloadObject(p.Contents[i])
				if err != nil {
					errc <- err
					return false
				}
				obj.ID = i
				// do something with the object
				fmt.Printf("chunk %d downloaded\n", obj.ID)
				handler.HandleChunk(obj)
			}
			return true
		})
		if err != nil {
			errc <- err
			return
		}

		handler.OnDone()
		fmt.Println("download done")
	}()

	return cntc, errc
}

func (dldr *s3Downloader) downloadObject(o *s3.Object) (*Chunk, error) {
	obj := new(Chunk)
	params := s3.GetObjectInput{
		Bucket: &dldr.bucket,
		Key:    o.Key,
	}

	if _, err := dldr.dloader.Download(obj, &params); err != nil {
		return nil, err
	}

	return obj, nil
}
