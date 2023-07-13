package drivers

import (
	"github.com/daqiancode/env"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIO() (*minio.Client, error) {
	return minio.New(env.Get("S3_ENDPOINT"), &minio.Options{
		Region: env.Get("S3_REGION"),
		Creds:  credentials.NewStaticV4(env.Get("S3_ACCESS_KEY"), env.Get("S3_SECRET_KEY"), env.Get("S3_SECRET_TOKEN")),
		Secure: env.GetBoolMust("S3_USE_SSL", true),
	})
}
