package stypes

import "strings"

type Language string

func (l Language) String() string {
	return strings.ToLower(string(l))
}

const (
	Zh Language = "zh"
	En Language = "en"
	Ar Language = "ar"
)
