package hashs_test

import (
	"testing"

	"github.com/scilive/scibase/utils/hashs"
	"github.com/stretchr/testify/assert"
)

func TestHashId(t *testing.T) {
	hashId := hashs.NewHashID("salt", 10)
	encoded := hashId.Encode(1)
	decoded, _ := hashId.Decode(encoded)
	assert.Equal(t, int64(1), decoded)
}
