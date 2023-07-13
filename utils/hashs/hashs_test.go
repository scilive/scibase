package hashs_test

import (
	"testing"

	"github.com/scilive/scibase/utils/hashs"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	assert.Equal(t, "e10adc3949ba59abbe56e057f20f883e", hashs.Md5("123456"))
	assert.Equal(t, "7c4a8d09ca3762af61e59520943dc26494f8941b", hashs.Sha1("123456"))
}
