package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateDynamoID(prefix string) (string, error) {
	buffer := make([]byte, 16)

	_, err := rand.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("error generating Dynamo ID: %w", err)
	}

	return prefix + hex.EncodeToString(buffer), nil
}
