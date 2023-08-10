package langs

import "strings"

type Country string

func (c Country) String() string {
	return strings.ToUpper(string(c))
}

const (
	CN Country = "zh"
	US Country = "en"
	AE Country = "ar"
)
