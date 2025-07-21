package tokenstore

import (
	"sync"
	"time"
)

// Store refresh token to create, track and revoke sessions
type RefreshTokenStore interface {
	Set(token string, userID int64, duration int)
	Get(token string) (userID int64, found bool)
	Delete(token string)
}

type memoryTokenStore struct {
	tokens tokenStore
	lock   sync.RWMutex
}

type tokenData struct {
	userID    int64
	expiresAt time.Time
}

type tokenStore = map[string]tokenData
