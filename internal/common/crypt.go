package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
)

func ComplicatedPasswordForPrepareSign(password []byte) []byte {
	salt := []byte{0x24, 0x4a, 0xf0, 0xbb, 0xc9, 0xa4, 0xe7, 0x56, 0x1f, 0x02, 0x1a, 0xb4}
	key := argon2.IDKey(password, salt, 1, 16*1024, 4, 32)
	return key
}

func CalculateSignKey(secret []byte, salt []byte) []byte {
	hash := sha256.New
	hkdf := hkdf.New(hash, secret, salt, nil)
	key := make([]byte, 32)
	io.ReadFull(hkdf, key)
	return key
}
func GetKeyFromString(key string) ([]byte, error) {
	return []byte(key), nil
}
func Sign(msg []byte, signKey []byte) ([]byte, error) {
	h := hmac.New(sha256.New, signKey)
	h.Write(msg)
	return h.Sum(nil), nil
}

func CheckHash(msg []byte, msgSign []byte, signKey []byte) bool {
	actualSign, _ := Sign(msg, signKey)
	return hmac.Equal(actualSign, msgSign)
}

var AES256IVSIZE int = 16

func EncryptAES(dataForSafe []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, fmt.Errorf("can't create chipher %w", err)
	}
	iv := make([]byte, AES256IVSIZE)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, fmt.Errorf("can't create iv %w", err)
	}
	safeData := make([]byte, len(dataForSafe))
	ofbStream := cipher.NewOFB(c, iv)
	ofbStream.XORKeyStream(safeData, dataForSafe)

	var cipherText []byte
	cipherText = append(cipherText, iv...)
	cipherText = append(cipherText, safeData...)

	return cipherText, nil
}

func DecryptAES(cipherText []byte, key []byte) ([]byte, error) {
	iv := cipherText[:AES256IVSIZE]
	c, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, fmt.Errorf("can't create chipher %v", err)
	}
	ofbStream := cipher.NewOFB(c, iv)
	temp := cipherText[AES256IVSIZE:]
	plaintext := make([]byte, len(cipherText)-AES256IVSIZE)
	ofbStream.XORKeyStream(plaintext, temp)
	return plaintext, nil
}
