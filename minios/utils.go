package minios

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"

	"github.com/daqiancode/env"
	"github.com/minio/minio-go/v7"
	"github.com/scilive/scibase/drivers"
	"github.com/scilive/scibase/utils/images"
	"github.com/scilive/scibase/utils/rands"
)

func Put(file io.Reader, category, fileName, contentType string, fileSize int64, threads uint) (string, error) {
	client, err := NewMiniosEnv()
	if err != nil {
		return "", err
	}
	return client.Put(file, category, fileName, contentType, fileSize, threads)
	// s3, err := drivers.NewMinIO()
	// if err != nil {
	// 	return "", err
	// }
	// s3Client := &MinIOClient{Client: s3, Bucket: env.Get("S3_BUCKET")}
	// path := filepath.Join(env.Get("S3_ROOT"), category, rands.RandomPath())
	// fullPath := path + filepath.Ext(fileName)
	// err = s3Client.Put(fullPath, file, fileSize, contentType, "", threads)
	// if err != nil {
	// 	return "", err
	// }
	// return fullPath, nil
}

func Get(key string) ([]byte, error) {
	client, err := NewMiniosEnv()
	if err != nil {
		return nil, err
	}
	return client.Get(key)
	// s3, err := drivers.NewMinIO()
	// if err != nil {
	// 	return nil, err
	// }
	// s3Client := &MinIOClient{Client: s3, Bucket: env.Get("S3_BUCKET")}
	// res, err := s3Client.Get(key)
	// if err != nil {
	// 	return nil, err
	// }
	// defer res.Close()
	// return io.ReadAll(res)
}

type Minios struct {
	Client *MinIOClient
	// Bucket  string
	RootDir string
}

func NewMiniosEnv() (*Minios, error) {
	s3, err := drivers.NewMinIO()
	if err != nil {
		return nil, err
	}
	return &Minios{
		Client:  &MinIOClient{Client: s3, Bucket: env.Get("S3_BUCKET")},
		RootDir: env.Get("S3_ROOT"),
	}, nil
}

func (s *Minios) Put(file io.Reader, category, ext, contentType string, fileSize int64, threads uint) (string, error) {
	path := filepath.Join(s.RootDir, category, rands.RandomPath())
	if !strings.Contains(ext, ".") {
		ext = "." + ext
	} else {
		ext = filepath.Ext(ext)
	}
	fullPath := path + ext
	err := s.Client.Put(fullPath, file, fileSize, contentType, "", threads)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}
func (s *Minios) Puts(files []io.Reader, category, ext, contentType string, fileSizes []int64, threads uint) ([]string, error) {
	var paths []string
	for i, file := range files {
		path, err := s.Put(file, category, ext, contentType, fileSizes[i], threads)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func (s *Minios) Get(key string) ([]byte, error) {
	res, err := s.Client.Get(key)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	return io.ReadAll(res)
}
func (s *Minios) Remove(key string) error {
	return s.Client.Remove(key)
}
func (s *Minios) Download(key, dst string) error {
	return s.Client.FGetObject(context.Background(), s.Client.Bucket, key, dst, minio.GetObjectOptions{})
}

type Rect struct {
	X, Y, W, H int
}
type PutImageResult struct {
	Raw     string
	Crop    string
	Resizes []string
}

func (s *Minios) PutImage(file io.Reader, category, fileName, contentType string, fileSize int64, crop Rect, resizes [][]int) (PutImageResult, error) {
	input := file
	var r PutImageResult
	var err error
	r.Raw, err = s.Put(input, category, fileName, contentType, fileSize, 0)
	if err != nil {
		return r, err
	}
	//crop
	if crop.H > 0 {
		bs, err := images.Crop(input, fileName, crop.X, crop.Y, crop.W, crop.H)
		if err != nil {
			return r, err
		}
		input = bytes.NewBuffer(bs)
	}
	//resize
	for _, v := range resizes {
		w, h := v[0], v[1]
		bs, err := images.Resize(input, fileName, w, h)
		if err != nil {
			return r, err
		}
		input = bytes.NewBuffer(bs)
	}
	return r, nil
}
