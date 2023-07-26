package images_test

import (
	"os"
	"testing"

	"github.com/daqiancode/env"
	"github.com/scilive/scibase/utils/images"
	"github.com/stretchr/testify/assert"
)

func TestResize(t *testing.T) {
	file, err := os.Open(env.GetPath("utils/images/test.jpg"))
	assert.Nil(t, err)
	bytes, err := images.Resize(file, "test.jpg", 200, 200)
	assert.Nil(t, err)
	err = os.WriteFile(env.GetPath("utils/images/test_reszied.jpg"), bytes, 0644)
	assert.Nil(t, err)
}

func TestCrop(t *testing.T) {
	file, err := os.Open(env.GetPath("utils/images/test.jpg"))
	assert.Nil(t, err)
	bytes, err := images.Crop(file, "test.jpg", 100, 100, 200, 200)
	assert.Nil(t, err)
	err = os.WriteFile(env.GetPath("utils/images/test_cropped.jpg"), bytes, 0644)
	assert.Nil(t, err)
}

func TestThumbnail(t *testing.T) {
	file, err := os.Open(env.GetPath("utils/images/test.jpg"))
	assert.Nil(t, err)
	bytes, err := images.Thumbnail(file, "test.jpg", 200, 200)
	assert.Nil(t, err)
	err = os.WriteFile(env.GetPath("utils/images/test_thumbnail.jpg"), bytes, 0644)
	assert.Nil(t, err)
}
