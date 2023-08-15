package stypes

import "strings"

type Country string

func (c Country) String() string {
	return strings.ToUpper(string(c))
}

const (
	CN Country = "CN"
	US Country = "EN"
	AE Country = "AE"
)
