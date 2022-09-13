package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func TokenHex(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error reading random token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// GetRandomDuration creates a random duration between `low` and `high` seconds (inclusive)
func GetRandomDuration(low int, high int) time.Duration {
	high *= 1000
	low *= 1000
	r := rand.Intn(high-low+1) + low
	return time.Duration(r) * time.Millisecond
}

// Remove removes element `e` from list `s`
func Remove[T comparable](s []T, e T) []T {
	lst := make([]T, 0)
	for _, i := range s {
		if i != e {
			lst = append(lst, i)
		}
	}
	return lst
}
