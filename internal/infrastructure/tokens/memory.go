// Refresh token store implemented to use in-memory

package tokens

import (
	"sync"
	"time"
)

type memoryStore struct {
	tokens tokenStore
	lock   sync.RWMutex
}

var (
	once     sync.Once
	instance Store
)

func GetStore() Store {
	once.Do(func() {
		instance = newMemoryStore()
	})
	return instance
}

func newMemoryStore() Store {
	return &memoryStore{
		tokens: make(tokenStore),
	}
}

func (store *memoryStore) Set(
	token string,
	userId int64,
	duration int,
) {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.tokens[token] = tokenData{
		userId:    userId,
		expiresAt: time.Now().Add(time.Duration(duration) * time.Second),
	}
}

func (store *memoryStore) Get(token string) (int64, bool) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	data, exists := store.tokens[token]
	if !exists || time.Now().After(data.expiresAt) {
		return 0, false
	}

	return data.userId, true
}

func (store *memoryStore) Delete(token string) {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.tokens, token)
}
