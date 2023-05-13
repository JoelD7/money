package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
)

var errHashMismatch = errors.New("token and hash mismatch")

func hashTokens(accessToken, refreshToken string) (string, string, error) {
	h := sha256.New()

	_, err := h.Write([]byte(accessToken))
	if err != nil {
		return "", "", err
	}

	hashedAccess := h.Sum(nil)
	h.Reset()

	_, err = h.Write([]byte(refreshToken))
	if err != nil {
		return "", "", err
	}

	hashedRefresh := h.Sum(nil)
	h.Reset()

	return fmt.Sprintf("%x", hashedAccess), fmt.Sprintf("%x", hashedRefresh), nil
}

func compareHashAndToken(hash, token string) error {
	h := sha256.New()

	_, err := h.Write([]byte(token))
	if err != nil {
		return err
	}

	hashedToken := h.Sum(nil)

	if strings.Compare(hash, fmt.Sprintf("%x", hashedToken)) != 0 {
		return errHashMismatch
	}

	return nil
}
