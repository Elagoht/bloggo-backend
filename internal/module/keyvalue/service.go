package keyvalue

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/keyvalue/models"
	"bloggo/internal/module/webhook"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/audit"
	auditmodels "bloggo/internal/module/audit/models"
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
	userId int64,
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
	err = service.repository.BulkUpsert(items)
	if err != nil {
		return err
	}

	// Log the audit event (using 0 as entity ID since it's a bulk operation)
	audit.LogAction(&userId, auditmodels.EntityKeyValue, 0, auditmodels.ActionUpdated)

	// Trigger webhook with updated key-value pairs
	keyValueMap := make(map[string]interface{})
	for _, item := range items {
		keyValueMap[item.Key] = item.Value
	}
	go func() { webhook.TriggerKeyValueUpdated(keyValueMap) }()

	return nil
}
