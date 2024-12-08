package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateSKU() string {
	timestamp := time.Now().UnixNano()

	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random bytes")
	}

	randomPart := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("SKU-%d-%s", timestamp, randomPart)
}
