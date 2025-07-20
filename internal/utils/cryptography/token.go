package cryptography

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Creates a unique random string for refresh token.
func GenerateUniqueId() (string, error) {
	uniqueId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uniqueId.String(), nil
}

// Creates a JWT with the given expiry and type.
func GenerateJWT(
	subject string,
	userId int64,
	roleId int64,
	secret string,
	duration int,
) (string, error) {

	claims := jwt.MapClaims{
		"sub": subject,
		"uid": userId,
		"rid": roleId,
		"exp": time.Now().Add(time.Duration(duration)).Unix(),
		"iat": time.Now().Unix(),
		"iss": "bloggo",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}
