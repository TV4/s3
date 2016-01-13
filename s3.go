package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// ChunkHandler declares the interface needed to receive the chunks downloaded
// by the Downloader.
type ChunkHandler interface {
	HandleChunk(*Chunk)
	OnDone()
}

// BucketConf represents the data needed to access an S3 bucket. Specifically,
// it contains the bucket name, region, and key/secret credentials
type BucketConf struct {
	Bucket, ID, Secret, Region string
}

func awsConfig(id, secret, region string) *aws.Config {
	return &aws.Config{
		Credentials: credentials.NewChainCredentials([]credentials.Provider{
			&credentials.StaticProvider{
				Value: credentials.Value{AccessKeyID: id, SecretAccessKey: secret},
			},
		}),
		Logger: aws.NewDefaultLogger(),
		Region: &region,
	}
}
