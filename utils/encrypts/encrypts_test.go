package encrypts_test

import (
	"testing"

	"github.com/scilive/scibase/utils/encrypts"
	"github.com/stretchr/testify/assert"
)

func TestRSA(t *testing.T) {
	crtypedPassword := encrypts.Bcrypt("123456", 10)
	assert.True(t, encrypts.BcryptCompare("123456", crtypedPassword))
}
