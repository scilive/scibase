package strs_test

import (
	"testing"

	"github.com/scilive/scibase/utils/strs"
	"github.com/stretchr/testify/assert"
)

func TestCaptitalize(t *testing.T) {
	r := strs.Capitalize("hello")
	assert.Equal(t, "Hello", r)
	r = strs.Capitalize("")
	assert.Equal(t, "", r)
}

func TestDecapitalize(t *testing.T) {
	r := strs.Decapitalize("HELLO")
	assert.Equal(t, "hELLO", r)
	r = strs.Decapitalize("")
	assert.Equal(t, "", r)
}
