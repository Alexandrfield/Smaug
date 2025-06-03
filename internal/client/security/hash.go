package security

import (
	"crypto/sha256"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func HashInputPassword(password []byte) []byte {
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		keyLength:   32,
	}
	salt := []byte{0xf4, 0x51, 0x78, 0x01, 0xbc, 0x00, 0xd6, 0xa2, 0xd1, 0xe9, 0xaa, 0x01, 0xd0, 0x00, 0xfa, 0x0cd}
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	return hash
}

func generateCryptoKey(login []byte, password []byte) []byte {
	h := sha256.New()
	h.Write([]byte{0x02, 0xfa, 0xa4, 0x55})
	h.Write(login)
	h.Write(password)
	return h.Sum(nil)
}
