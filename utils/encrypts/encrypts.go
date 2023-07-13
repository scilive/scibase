package encrypts

import (
	"golang.org/x/crypto/bcrypt"
)

func Bcrypt(password string, cost int) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func BcryptCompare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
