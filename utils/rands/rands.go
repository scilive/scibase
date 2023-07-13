package rands

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var numbers = []byte("0123456789")

func UUID(length ...int) string {
	l := 32
	if len(length) > 0 {
		l = length[0]
	}
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomInts(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}

// RandomPath  returns Jc/Sp/hmDoWw5BTISBCHhCzwXj.jpg
func RandomPath(suffix string) string {
	p := UUID(4)
	i := strings.Index(suffix, ".")
	if i == -1 {
		i = 0
	} else {
		i += 1
	}
	return fmt.Sprintf("%s/%s/%s.%s", p[:2], p[2:], UUID(20), suffix[i:])
}
