// Refresh token store implemented to use in-memory

package tokenstore

import (
	"sync"
	"time"
)

var (
	once     sync.Once
	instance RefreshTokenStore
)

func GetRefreshTokenStore() RefreshTokenStore {
	once.Do(func() {
		instance = newMemoryStore()
	})
	return instance
}

func newMemoryStore() RefreshTokenStore {
	return &memoryTokenStore{
		tokens: make(tokenStore),
	}
}

func (store *memoryTokenStore) Set(
	token string,
	userID int64,
	duration int,
) {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.tokens[token] = tokenData{
		userID:    userID,
		expiresAt: time.Now().Add(time.Duration(duration) * time.Second),
	}
}

func (store *memoryTokenStore) Get(token string) (int64, bool) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	data, exists := store.tokens[token]
	if !exists || time.Now().After(data.expiresAt) {
		return 0, false
	}

	return data.userID, true
}

func (store *memoryTokenStore) Delete(token string) {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.tokens, token)
}
