package stypes_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scilive/scibase/stypes"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	cases := []string{
		"2023-08",
		"2023-08-24",
		"2023-08-24 12:00",
		"2023-08-24 12:00:00",
		"12:00:00",
		"12:00",
		"2023-08-24T12:00+04:00",
		"2023-08-24+04:00",
		"2023-08-24 12:00+04:00",
		"2023-08-24 12:00:00+04:00",
		"2023-08-24T12:00:00+04:00",
	}
	var r stypes.Time
	for _, c := range cases {
		c = fmt.Sprintf("\"%s\"", c)
		err := json.Unmarshal([]byte(c), &r)
		assert.Nil(t, err)
		_, err = json.Marshal(r)
		assert.Nil(t, err)
	}
}

func TestDate(t *testing.T) {
	str := `"2023-08-24"`
	var dt stypes.Date
	err := json.Unmarshal([]byte(str), &dt)
	assert.Nil(t, err)
	bs, err := json.Marshal(dt)
	assert.Nil(t, err)
	assert.Equal(t, str, string(bs))
}
