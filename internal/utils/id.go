package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomID generates a random string ID using crypto/rand
func GenerateRandomID() string {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a simple ID if rand fails (unlikely)
		return "fallback-id"
	}
	return hex.EncodeToString(bytes)
}
