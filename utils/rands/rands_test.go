package rands_test

import (
	"testing"

	"github.com/scilive/scibase/utils/rands"
	"github.com/stretchr/testify/assert"
)

func TestRandomPath(t *testing.T) {
	a := rands.RandomPath("jpg")
	assert.Equal(t, 30, len(a))
	b := rands.RandomPath(".jpg")
	assert.Equal(t, 30, len(b))
}
