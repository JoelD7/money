package hash

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"strings"
)

var (
	h hash.Hash

	ErrHashMismatch = errors.New("token and hash mismatch")
)

func init() {
	h = sha256.New()
}

// Apply applies the hash function to the token
func Apply(token string) (string, error) {
	defer h.Reset()

	_, err := h.Write([]byte(token))
	if err != nil {
		return "", err
	}

	hashedToken := h.Sum(nil)

	return fmt.Sprintf("%x", hashedToken), nil
}

// CompareWithToken determines if the passed hash value is the hashed value of the token
func CompareWithToken(hashValue, token string) error {
	defer h.Reset()

	_, err := h.Write([]byte(token))
	if err != nil {
		return err
	}

	hashedToken := h.Sum(nil)

	if strings.Compare(hashValue, fmt.Sprintf("%x", hashedToken)) != 0 {
		return ErrHashMismatch
	}

	return nil
}
