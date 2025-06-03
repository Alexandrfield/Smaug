package security

import (
	"testing"

	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestNewSafetyVault(t *testing.T) {
	key := []byte{0xff}
	p := NewSafetyVault(key, &common.FakeLogger{})
	if p == nil {
		t.Errorf("expected not nil result")
	}
}

func TestCreateKeys(t *testing.T) {
	key := []byte{0xff}
	p := NewSafetyVault(key, &common.FakeLogger{})
	if p == nil {
		t.Errorf("expected not nil result")
	}
	password := []byte{0x45, 0x55, 0xf1, 0x6d}
	plainText := []byte{0x00, 0x00, 0x00, 0x00}
	actualNewKey, _ := p.CreateKeys(password)
	encrypted := p.EncryptData(plainText, actualNewKey)
	decrypted := p.DecryptData(encrypted, actualNewKey)
	assert.ElementsMatch(t, plainText, decrypted)
}

func TestCcheckSign(t *testing.T) {
	key := []byte{0xff}
	p := NewSafetyVault(key, &common.FakeLogger{})
	if p == nil {
		t.Errorf("expected not nil result")
	}
	password := []byte{0x45, 0x55, 0xf1, 0x6d}
	dataxt := []byte{0x00, 0x00, 0x00, 0x00, 0xaa, 0xaa, 0xaa, 0xaa}
	temporyKey := []byte{0xad, 0xf0, 0xbe, 0x00, 0xad, 0xaa, 0xa1, 0x22}
	_, actualSignKeyComplicated := p.CreateKeys(password)
	signKey := common.CalculateSignKey(actualSignKeyComplicated, temporyKey)
	sign, _ := common.Sign(dataxt, signKey)
	if !p.CheckSign(actualSignKeyComplicated, temporyKey, dataxt, sign) {
		t.Errorf("check failed")
	}
}
