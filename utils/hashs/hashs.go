package hashs

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Md5(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
func Sha1(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
