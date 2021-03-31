package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

// NewUUID .
func NewUUID() (uuid.UUID, error) {
	return uuid.NewUUID()
}

// NewSalt .
func NewSalt() string {
	var b [16]byte
	rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
