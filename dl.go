package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Downloader struct {
	Cntc      chan int
	Errc      chan error
	bucket    string
	s3Client  *s3.S3
	dloader   *s3manager.Downloader
	maxchunks int
}

func NewS3Downloader(bucket, id, secret, region string, maxchunks int) s3Downloader {
	s3Client := s3.New(session.New(), awsConfig(id, secret, region))
	dloader := &s3manager.Downloader{
		PartSize:    1024 * 1024 * 5,
		Concurrency: 5,
		S3:          s3Client,
	}

	return s3Downloader{
		bucket:    bucket,
		s3Client:  s3Client,
		dloader:   dloader,
		maxchunks: maxchunks,
		Cntc:      make(chan int),
		Errc:      make(chan error),
	}
}

func (dldr *s3Downloader) DownloadChunks(path string, handler ChunkHandler) {
	objPath := &s3.ListObjectsInput{Bucket: &dldr.bucket, Prefix: &path}

	go func() {
		// Iterate over objects located in objPath
		err := dldr.s3Client.ListObjectsPages(objPath, func(p *s3.ListObjectsOutput, lastPage bool) bool {
			n := dldr.maxchunks
			if dldr.maxchunks <= 0 || dldr.maxchunks > len(p.Contents) {
				n = len(p.Contents)
			}
			dldr.Cntc <- n
			for i := 0; i < n; i++ {
				obj, err := dldr.downloadObject(p.Contents[i])
				if err != nil {
					fmt.Println(err)
					dldr.Errc <- err
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
			fmt.Println("XXXX", err)
			dldr.Errc <- err
			return
		}

		handler.Stop()
		fmt.Println("download done")
	}()

	return
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
