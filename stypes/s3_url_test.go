package stypes_test

import (
	"encoding/json"
	"testing"

	"github.com/scilive/scibase/stypes"
	"github.com/stretchr/testify/assert"
)

func TestS3Url(t *testing.T) {
	key := "a/b.txt"
	s3url := stypes.S3Url(key)
	assert.Equal(t, "https://sci.live/a/b.txt", s3url.Url())
	bs, err := json.Marshal(s3url)
	assert.Nil(t, err)
	assert.Equal(t, "\"https://sci.live/a/b.txt\"", string(bs))
	err = json.Unmarshal(bs, &s3url)
	assert.Nil(t, err)
}

func TestOtherTypes(t *testing.T) {
	lang := stypes.Language("en")
	assert.Equal(t, "en", lang.String())
	country := stypes.Country("US")
	assert.Equal(t, "US", country.String())
}

func TestTag(t *testing.T) {
	tags := stypes.Tags{{Name: "tag1"}}
	value, err := tags.Value()
	assert.Nil(t, err)
	var tags2 stypes.Tags
	err = tags2.Scan(value)
	assert.Nil(t, err)
	assert.Equal(t, tags, tags2)

}
