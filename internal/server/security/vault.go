package security

import (
	"crypto/rand"
	"io"

	"github.com/Alexandrfield/Smaug/internal/common"
)

func NewSafetyVault(key []byte, logger common.Logger) *SafetyVault {
	for i := len(key); i < 32; i++ { // expand AES key for 256 bit
		key = append(key, 0)
	}
	return &SafetyVault{masterHidenKey: key, logger: logger}
}

type SafetyVault struct {
	masterHidenKey []byte
	logger         common.Logger
}

func (vault *SafetyVault) CreateKeys(password []byte) ([]byte, []byte) {
	key := make([]byte, 32)
	io.ReadFull(rand.Reader, key)
	newKey, err := common.EncryptAES(key, vault.masterHidenKey)
	if err != nil {
		vault.logger.Debugf("encrypt error:%s", err)
	}
	signKeyComplicated := common.ComplicatedPasswordForPrepareSign(password)
	return newKey, signKeyComplicated
}
func (vault *SafetyVault) EncryptData(plainDta []byte, key []byte) []byte {
	realKey, err := common.DecryptAES(key, vault.masterHidenKey)
	if err != nil {
		vault.logger.Debugf("decrypt error:%s", err)
	}
	encryptedData, _ := common.EncryptAES(plainDta, realKey)
	return encryptedData
}
func (vault *SafetyVault) DecryptData(encryptedData []byte, key []byte) []byte {
	realKey, _ := common.DecryptAES(key, vault.masterHidenKey)
	decryptedData, _ := common.DecryptAES(encryptedData, realKey)
	return decryptedData
}
func (vault *SafetyVault) CheckSign(signKeyComplicated []byte, tempKey []byte, data []byte, actualSign []byte) bool {
	signKey := common.CalculateSignKey(signKeyComplicated, tempKey)
	return common.CheckHash(data, actualSign, signKey)
}
