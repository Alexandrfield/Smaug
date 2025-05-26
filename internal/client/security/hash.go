package security

import "golang.org/x/crypto/bcrypt"

func HashInputPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes)
}
