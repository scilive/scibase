package minios

import (
	"io"
	"path/filepath"

	"github.com/daqiancode/env"
	"github.com/scilive/scibase/drivers"
	"github.com/scilive/scibase/utils/rands"
)

func Put(file io.Reader, category, fileName, contentType string, fileSize int64, threads uint) (string, error) {
	s3, err := drivers.NewMinIO()
	if err != nil {
		return "", err
	}
	s3Client := &MinIOClient{Client: s3, Bucket: env.Get("S3_BUCKET")}
	path := filepath.Join(env.Get("S3_ROOT"), category, rands.RandomPath())
	fullPath := path + filepath.Ext(fileName)
	err = s3Client.Put(fullPath, file, fileSize, contentType, "", threads)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}

func Get(key string) ([]byte, error) {
	s3, err := drivers.NewMinIO()
	if err != nil {
		return nil, err
	}
	s3Client := &MinIOClient{Client: s3, Bucket: env.Get("S3_BUCKET")}
	res, err := s3Client.Get(key)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	return io.ReadAll(res)
}
