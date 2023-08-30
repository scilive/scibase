package rands

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/speps/go-hashids"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var lettersLower = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
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

func UUIDLower(length ...int) string {
	l := 32
	if len(length) > 0 {
		l = length[0]
	}
	b := make([]byte, l)
	for i := range b {
		b[i] = lettersLower[rand.Intn(len(lettersLower))]
	}
	return string(b)
}

// RandomPath  returns /Jc/Sp/hmDoWw5BTISBCHhCzwXj
func RandomPath() string {
	p := UUID(4)
	return fmt.Sprintf("/%s/%s/%s", p[:2], p[2:], UUID(20))
}
func RandomDatePath() string {
	return fmt.Sprintf("/%s/%s", DayHash(), UUID(10))
}

func IdPath(id int64) string {
	return fmt.Sprintf("/%02x/%s", id%256, HashId8(id))
}

func DayHash() string {
	now := time.Now()
	cur, _ := strconv.ParseInt(fmt.Sprintf("%d%02d%02d", now.Year(), now.Month(), now.Day()), 10, 64)
	return HashId4(cur)
}

var salt = "__sci_live@2023__"

func HashId4(id int64) string {
	generator := NewHashIdGenerator(salt, 4)
	return generator.Encode(id)
}

func DeHashId4(id string) (int64, error) {
	generator := NewHashIdGenerator(salt, 4)
	return generator.Decode(id)
}

func HashId8(id int64) string {
	generator := NewHashIdGenerator(salt, 8)
	return generator.Encode(id)
}

func DeHashId8(id string) (int64, error) {
	generator := NewHashIdGenerator(salt, 8)
	return generator.Decode(id)
}

type HashIdGenerator struct {
	minLength int
	salt      string
	h         *hashids.HashID
}

// NewHashIdGenerator create a new HashIdGenerator
func NewHashIdGenerator(salt string, minLength int) *HashIdGenerator {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, _ := hashids.NewWithData(hd)
	return &HashIdGenerator{
		minLength: minLength,
		salt:      salt,
		h:         h,
	}
}

// Encoder of hashId
func (s HashIdGenerator) Encode(id int64) string {
	cur, _ := s.h.EncodeInt64([]int64{id})
	return cur
}

// Decoder of hashId
func (s HashIdGenerator) Decode(id string) (int64, error) {
	cur, err := s.h.DecodeInt64WithError(id)
	if err != nil {
		return 0, err
	}
	return cur[0], nil
}
