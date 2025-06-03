package security

import (
	"bytes"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "qwerftgh"
	b1 := HashPassword([]byte(password))
	b2 := HashPassword([]byte(password))
	if !bytes.Equal(b1, b2) {
		t.Errorf("Now equal. b1:%x; b2:%x", b1, b2)
		return
	}
}
