package utils

import (
	"crypto/rand"
	"fmt"
)

func GenerateUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
