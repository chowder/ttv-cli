package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func TokenHex(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error reading random token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
