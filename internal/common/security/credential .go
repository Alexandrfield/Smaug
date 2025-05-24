package security

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/Alexandrfield/Smaug/internal/common"
)

type Credential struct {
	logger common.Logger
	login  string
	secret []byte
}

func GetNewCredential(login string, password string, temporyKey []byte, logger common.Logger) *Credential {
	h := hmac.New(sha256.New, temporyKey)
	h.Write([]byte(password + login))

	t := Credential{logger: logger, secret: h.Sum(nil), login: login}
	return &t
}

func (cred *Credential) GetLogin() string {
	return cred.login
}
func (cred *Credential) EncryptData(data []byte) []byte {
	return data
}
func (cred *Credential) DecryptData(data []byte) []byte {
	return data
}
