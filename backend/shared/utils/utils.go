package utils

import (
	"encoding/json"
	"math/rand"
	"time"
)

// GenerateDynamoID generates a hex-based random unique ID with the given prefix
func GenerateDynamoID(prefix string) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 20)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return prefix + string(b)
}

// GetJsonString returns the json string representation of a given object
func GetJsonString(object interface{}) (string, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
