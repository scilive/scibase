package minios

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/scilive/scibase/logs"
)

type S3Provider string

const (
	S3ProviderAmazon     S3Provider = "AMAZON"
	S3ProviderAzure      S3Provider = "AZURE"
	S3ProviderCloudflare S3Provider = "CLOUDFLARE"
	S3ProviderHuawei     S3Provider = "HUAWEI"
)

type MinIOClient struct {
	*minio.Client
	Bucket   string
	Provider S3Provider
}

func NewMinIO(client *minio.Client, bucket string, provider S3Provider) *MinIOClient {
	return &MinIOClient{
		Client:   client,
		Bucket:   bucket,
		Provider: provider,
	}
}

func (s *MinIOClient) Put(key string, file io.Reader, fileSize int64, contentType, filename string) error {
	key = strings.TrimLeft(key, "/")
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}
	if filename != "" {
		opts.ContentDisposition = fmt.Sprintf("attachment; filename=\"%s\"", filename)
	}
	_, err := s.PutObject(context.Background(), s.Bucket, key, file, fileSize, opts)
	return err
}

func (s *MinIOClient) Get(key string) (*minio.Object, error) {
	key = strings.TrimLeft(key, "/")
	return s.Client.GetObject(context.Background(), s.Bucket, key, minio.GetObjectOptions{})
}

func (s *MinIOClient) Remove(key string) error {
	key = strings.TrimLeft(key, "/")
	return s.Client.RemoveObject(context.Background(), s.Bucket, key, minio.RemoveObjectOptions{})
}

func (s *MinIOClient) Url(key string) string {
	schema := s.Client.EndpointURL().Scheme
	endpoint := s.EndpointURL().Host
	if s.Provider == S3ProviderHuawei {
		return fmt.Sprintf("%s://%s.%s/%s", schema, s.Bucket, endpoint, key)
	}
	if s.Provider == S3ProviderAzure {
		return fmt.Sprintf("%s://%s.blob.%s/%s", schema, s.Bucket, endpoint, key)
	}
	if s.Provider == S3ProviderAmazon {
		return fmt.Sprintf("%s://%s.s3.%s/%s", schema, s.Bucket, endpoint, key)
	}
	if s.Provider == S3ProviderCloudflare {
		return fmt.Sprintf("%s://%s/%s", schema, endpoint, key)
	}
	return fmt.Sprintf("%s://%s/%s", schema, endpoint, key)
}

func CreateBucket(client *minio.Client, bucket string, publicRead bool) error {
	err := client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(context.Background(), bucket)
		if errBucketExists == nil && exists {
			logs.Log.Info().Str("bucket", bucket).Msg("bucket exists")
		} else {
			logs.Log.Error().Err(err).Str("bucket", bucket).Msg("create bucket")
			return err
		}
		if publicRead {
			//set bucket policy to public read
			client.SetBucketPolicy(context.Background(), bucket, `{"Version":"2012-10-17","Statement":[{"Sid":"PublicReadGetObject","Effect":"Allow","Principal":"*","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::`+bucket+`/*"]}]}`)
		}
	} else {
		logs.Log.Info().Str("bucket", bucket).Msg("create bucket success")
	}
	return nil
}
