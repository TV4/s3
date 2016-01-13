package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type ChunkHandler interface {
	HandleChunk(*Chunk)
	Stop()
}

type Config struct {
	Bucket, ID, Secret, Region string
	Nchunks                    int
}

func awsConfig(id, secret, region string) *aws.Config {
	return &aws.Config{
		Credentials: credentials.NewChainCredentials([]credentials.Provider{
			&credentials.StaticProvider{credentials.Value{AccessKeyID: id, SecretAccessKey: secret}},
		}),
		Logger: aws.NewDefaultLogger(),
		Region: &region,
	}
}
