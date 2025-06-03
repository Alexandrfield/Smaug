package security

import (
	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func HashPassword(password []byte) []byte {
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		keyLength:   32,
	}
	salt := []byte{0x11, 0x01, 0xdd, 0x01, 0x31, 0x45, 0xfa, 0x0ed, 0x11, 0x01, 0x1f, 0xe1, 0x00, 0x0e, 0x31, 0x20}
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	return hash
}
