package maths_test

import (
	"testing"

	"github.com/scilive/scibase/utils/maths"
	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	r := maths.Max(values...)
	assert.Equal(t, r, 5)
}

func TestMin(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	r := maths.Min(values...)
	assert.Equal(t, r, 1)
}
