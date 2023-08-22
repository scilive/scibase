package minios

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/scilive/scibase/logs"
	"github.com/scilive/scibase/utils/rands"
)

type MinIOClient struct {
	*minio.Client
	Bucket string
}

func NewMinIO(client *minio.Client, bucket string) *MinIOClient {
	return &MinIOClient{
		Client: client,
		Bucket: bucket,
	}
}

func (s *MinIOClient) Save(rootDir string, file io.Reader, fileSize int64, contentType, filename string, numThreads uint) (string, error) {
	key := filepath.Join(strings.TrimLeft(rootDir, "/"), rands.RandomPath())
	key += filepath.Ext(filename)
	err := s.Put(key, file, fileSize, contentType, filename, numThreads)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *MinIOClient) Put(key string, file io.Reader, fileSize int64, contentType, filename string, numThreads uint) error {
	key = strings.TrimLeft(key, "/")
	var opts minio.PutObjectOptions
	if contentType != "" {
		opts.ContentType = contentType
	}
	if filename != "" {
		opts.ContentDisposition = fmt.Sprintf("attachment; filename=\"%s\"", filename)
	}
	if numThreads > 1 {
		opts.NumThreads = numThreads
		opts.ConcurrentStreamParts = true
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
