package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Download starts the downloading of a resource residing at path in the
// bucket given by the configuration.
func Download(c BucketConf, path string, handler ObjectHandler) (<-chan int, <-chan error) {
	d := NewDownloader(c)
	return d.Download(path, handler)
}

// Downloader exposes two methods to download objects from an S3 bucket.
// Download takes as parameters a path to a resource located in the bucket,
// and an ObjectHandler to handle each object of the resource as it is
// downloaded. DownloadObjects is for testing purposes and offers a way to
// limit the number of objects downloaded. This can be convenient when dealing
// with resources containing large numbers of objects.
type Downloader interface {
	Download(path string, handler ObjectHandler) (<-chan int, <-chan error)
	DownloadObjects(path string, handler ObjectHandler, nobj, startobj int) (<-chan int, <-chan error)
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

func (d s3Downloader) Download(path string, handler ObjectHandler) (<-chan int, <-chan error) {
	return d.DownloadObjects(path, handler, 0, 0)
}

func (d s3Downloader) DownloadObjects(path string, handler ObjectHandler, nobj, startobj int) (<-chan int, <-chan error) {
	objPath := &s3.ListObjectsInput{Bucket: &d.bucket, Prefix: &path}
	cntc, errc := make(chan int), make(chan error)

	go func() {
		// Iterate over objects located in objPath
		err := d.s3Client.ListObjectsPages(objPath, func(p *s3.ListObjectsOutput, lastPage bool) bool {
			if nobj <= 0 || nobj > len(p.Contents)-startobj {
				nobj = len(p.Contents) - startobj
			}
			cntc <- nobj

			var (
				obj *Object
				err error
			)
			for i := startobj; i < startobj+nobj; i++ {
				obj, err = d.downloadObject(p.Contents[i])
				if err != nil {
					errc <- err
					return false
				}
				obj.ID = i
				handler.HandleObject(obj)
			}
			return true
		})
		if err != nil {
			errc <- err
			return
		}

		handler.OnDone()
	}()

	return cntc, errc
}

func (d s3Downloader) downloadObject(o *s3.Object) (*Object, error) {
	obj := new(Object)
	params := s3.GetObjectInput{
		Bucket: &d.bucket,
		Key:    o.Key,
	}

	if _, err := d.dloader.Download(obj, &params); err != nil {
		return nil, err
	}

	return obj, nil
}
