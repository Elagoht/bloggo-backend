package keyvalues

import (
	"bloggo/internal/module/api/keyvalues/models"
)

type KeyValuesAPIService struct {
	repository KeyValuesAPIRepository
}

func NewKeyValuesAPIService(repository KeyValuesAPIRepository) KeyValuesAPIService {
	return KeyValuesAPIService{repository}
}

func (service *KeyValuesAPIService) GetKeyValues(key, starting string) ([]models.APIKeyValue, error) {
	// If key parameter is provided, get specific key
	if key != "" {
		return service.repository.GetKeyValueByKey(key)
	}

	// If starting parameter is provided, get keys starting with prefix
	if starting != "" {
		return service.repository.GetKeyValuesStartingWith(starting)
	}

	// Otherwise, get all key-values
	return service.repository.GetAllKeyValues()
}
