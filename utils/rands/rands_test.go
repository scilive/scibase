package rands_test

import (
	"testing"

	"github.com/scilive/scibase/utils/rands"
	"github.com/stretchr/testify/assert"
)

func TestRandomPath(t *testing.T) {
	a := rands.RandomPath()
	assert.Equal(t, 26, len(a))

}

func TestUUID(t *testing.T) {
	a := rands.UUID()
	assert.Equal(t, 32, len(a))
}

func TestUUIDLower(t *testing.T) {
	a := rands.UUIDLower()
	assert.Equal(t, 32, len(a))
}
func TestRandomInts(t *testing.T) {
	a := rands.RandomInts(6)
	assert.Equal(t, 6, len(a))
}
