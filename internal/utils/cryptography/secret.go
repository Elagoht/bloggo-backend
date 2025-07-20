package cryptography

import (
	"crypto/rand"
	"encoding/base64"
)

// Cenerates a secure random 32-byte (256-bit) secret for HS256 JWT.
func GenerateRandomHS256Secret() (string, error) {
	byteArray := make([]byte, 24)
	_, err := rand.Read(byteArray)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(byteArray), nil
}
