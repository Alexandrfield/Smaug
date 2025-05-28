package security

import "golang.org/x/crypto/argon2"

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func HashInputPassword(password []byte) []byte {
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
	salt := []byte{0xf4, 0x51, 0x78, 0x01, 0xbc, 0x00, 0xd6, 0xa2, 0xd1, 0xe9, 0xaa, 0x01, 0xd0, 0x00, 0xfa, 0x0cd}
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	return hash
}
