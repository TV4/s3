package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Download starts the downloading of a resource residing at path in the
// bucket given by the configuration.
func Download(c BucketConf, path string, handler ChunkHandler) (<-chan int, <-chan error) {
	d := NewDownloader(c)
	return d.Download(path, handler)
}

// Downloader exposes two methods to download objects from an S3 bucket.
// Download takes as parameters a path to a resource located in the bucket,
// and a ChunkHandler to handle each chunk of the resources as it is
// downloaded. DownloadChunks is for testing purposes and offers a way to
// limit the number of chunks downloaded. This can be convenient when dealing
// with resources split into large numbers of chunks.
type Downloader interface {
	Download(path string, handler ChunkHandler) (<-chan int, <-chan error)
	DownloadChunks(path string, handler ChunkHandler, chunks int) (<-chan int, <-chan error)
}

type s3Downloader struct {
	bucket   string
	s3Client *s3.S3
	dloader  *s3manager.Downloader
}

// NewDownloader creates and initializes a Downloader for the specified bucket
// in the specified region.
func NewDownloader(c BucketConf) Downloader {
	return newS3Downloader(c.Bucket, c.ID, c.Secret, c.Region)
}

func newS3Downloader(bucket, id, secret, region string) s3Downloader {
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

func (d s3Downloader) Download(path string, handler ChunkHandler) (<-chan int, <-chan error) {
	return d.DownloadChunks(path, handler, 0)
}

func (d s3Downloader) DownloadChunks(path string, handler ChunkHandler, chunks int) (<-chan int, <-chan error) {
	objPath := &s3.ListObjectsInput{Bucket: &d.bucket, Prefix: &path}
	cntc, errc := make(chan int), make(chan error)

	go func() {
		// Iterate over objects located in objPath
		err := d.s3Client.ListObjectsPages(objPath, func(p *s3.ListObjectsOutput, lastPage bool) bool {
			if chunks <= 0 || chunks > len(p.Contents) {
				chunks = len(p.Contents)
			}
			cntc <- chunks
			for i := 0; i < chunks; i++ {
				obj, err := d.downloadObject(p.Contents[i])
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

func (d s3Downloader) downloadObject(o *s3.Object) (*Chunk, error) {
	obj := new(Chunk)
	params := s3.GetObjectInput{
		Bucket: &d.bucket,
		Key:    o.Key,
	}

	if _, err := d.dloader.Download(obj, &params); err != nil {
		return nil, err
	}

	return obj, nil
}
