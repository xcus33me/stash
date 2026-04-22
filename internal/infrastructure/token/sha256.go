package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type SHA256Generator struct{}

func (SHA256Generator) Generate() (raw, hash string, err error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}

	raw = hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(sum[:])

	return
}

func (SHA256Generator) Verify(raw, hash string) bool {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:]) == hash
}
