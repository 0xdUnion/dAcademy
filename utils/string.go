package utils

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
)

func RandomSecureString() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func IsSha256Hex(s string) bool {
	var sha256Regex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	if !sha256Regex.MatchString(s) {
		return false
	} else {
		return true
	}
}
