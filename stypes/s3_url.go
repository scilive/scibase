package stypes

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/daqiancode/env"
)

type S3Url string

func S3KeyToURL(key, provider, bucketHost, bucket string) string {
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

func (s S3Url) MarshalJSON() ([]byte, error) {
	provider := env.Get("S3_PROVIDER")
	bucketHost := env.Get("S3_BUCKET_HOST")
	bucket := env.Get("S3_BUCKET")
	return json.Marshal(S3KeyToURL(string(s), provider, bucketHost, bucket))
}

func (s *S3Url) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, s)
}
