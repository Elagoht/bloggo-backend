package keyvalue

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/keyvalue/models"
	"bloggo/internal/utils/apierrors"
)

type KeyValueService struct {
	repository  KeyValueRepository
	permissions permissions.Store
}

func NewKeyValueService(repository KeyValueRepository, permissions permissions.Store) KeyValueService {
	return KeyValueService{
		repository,
		permissions,
	}
}

func (service *KeyValueService) GetAll(userRoleId int64) ([]models.KeyValue, error) {
	// Check if user has permission to manage key-values
	hasPermission := service.permissions.HasPermission(userRoleId, "keyvalue:manage")
	if !hasPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetAll()
}

func (service *KeyValueService) BulkUpsert(
	items []models.RequestKeyValueUpsert,
	userRoleId int64,
) error {
	// Check if user has permission to manage key-values
	hasPermission := service.permissions.HasPermission(userRoleId, "keyvalue:manage")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	// First delete all existing entries
	err := service.repository.DeleteAll()
	if err != nil {
		return err
	}

	// Then insert/update all items
	return service.repository.BulkUpsert(items)
}
