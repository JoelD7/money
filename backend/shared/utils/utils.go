package utils

import (
	"encoding/json"
	"math/rand"
	"time"
)

func GenerateDynamoID(prefix string) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 20)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return prefix + string(b)
}

func GetJsonString(object interface{}) (string, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
