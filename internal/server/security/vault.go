package security

import (
	"crypto/rand"
	"io"

	"github.com/Alexandrfield/Smaug/internal/common"
)

func NewSafetyVault(key []byte) *SafetyVault {
	return &SafetyVault{masterHidenKey: key}
}

type SafetyVault struct {
	masterHidenKey []byte
}

func (vault *SafetyVault) CreateKeys(password []byte) ([]byte, []byte) {
	key := make([]byte, 32)
	io.ReadFull(rand.Reader, key)
	newKey, _ := common.EncryptAES(key, vault.masterHidenKey)
	signKeyComplicated := common.ComplicatedPasswordForPrepareSign(password)
	return newKey, signKeyComplicated
}
func (vault *SafetyVault) EncryptData(plainDta []byte, key []byte) string {
	realKey, _ := common.DecryptAES(key, vault.masterHidenKey)
	encryptedData, _ := common.EncryptAES(plainDta, realKey)
	return string(encryptedData)
}
func (vault *SafetyVault) DecryptData(encryptedData []byte, key []byte) string {
	realKey, _ := common.DecryptAES(key, vault.masterHidenKey)
	decryptedData, _ := common.DecryptAES(encryptedData, realKey)
	return string(decryptedData)
}
func (vault *SafetyVault) CheckSign(signKeyComplicated []byte, tempKey []byte, data []byte, actualSign []byte) bool {
	signKey := common.CalculateSignKey(signKeyComplicated, tempKey)
	return common.CheckHash(data, actualSign, signKey)
}
