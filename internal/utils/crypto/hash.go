package cryptography

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassphrase hashes a passphrase using bcrypt.
func HashPassphrase(passphrase string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassphrase compares a hashed passphrase with a plain passphrase.
func ComparePassphrase(hashedPassphrase, passphrase string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassphrase), []byte(passphrase))
	return err == nil
}
