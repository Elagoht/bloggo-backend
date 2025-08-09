package tokens

import (
	"time"
)

// Store refresh token to create, track and revoke sessions
type Store interface {
	Set(token string, userId int64, duration int)
	Get(token string) (userId int64, found bool)
	Delete(token string)
}

type tokenData struct {
	userId    int64
	expiresAt time.Time
}

type tokenStore = map[string]tokenData
