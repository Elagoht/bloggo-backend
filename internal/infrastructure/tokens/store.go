package tokens

import (
	"time"
)

// Store refresh token to create, track and revoke sessions
type Store interface {
	Set(token string, userID int64, duration int)
	Get(token string) (userID int64, found bool)
	Delete(token string)
}

type tokenData struct {
	userID    int64
	expiresAt time.Time
}

type tokenStore = map[string]tokenData
