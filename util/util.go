package util

import (
	"crypto/rand"
	"fmt"
)

func GenerateToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
