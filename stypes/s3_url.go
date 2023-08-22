package stypes

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/daqiancode/env"
	"github.com/scilive/scibase/logs"
)

type S3Url string

func S3KeyToURL(key, provider, bucketHost, bucket string) string {
	if key == "" {
		return ""
	}
	if bucketHost == "" {
		logs.Log.Error().Msg("S3_BUCKET_HOST is empty")
	}
	if strings.HasPrefix(key, "http") {
		return key
	}
	provider = strings.ToLower(provider)
	if provider == "huaweiobs" || provider == "huawei" {
		return bucketHost + filepath.Join("/", bucket, key)
	}
	if provider == "storj" {
		return bucketHost + filepath.Join("/", key)
	}
	return bucketHost + filepath.Join("/", bucket, key)
}
func (s S3Url) Url() string {
	provider := env.Get("S3_PROVIDER")
	bucketHost := env.Get("S3_BUCKET_HOST")
	bucket := env.Get("S3_BUCKET")
	return S3KeyToURL(string(s), provider, bucketHost, bucket)
}

func (s S3Url) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Url())
}

func (s *S3Url) UnmarshalJSON(data []byte) error {
	var r string
	err := json.Unmarshal(data, &r)
	if err != nil {
		return err
	}
	*s = S3Url(r)
	return nil
}
