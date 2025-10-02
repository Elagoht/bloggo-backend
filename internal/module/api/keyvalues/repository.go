package keyvalues

import (
	"bloggo/internal/module/api/keyvalues/models"
	"database/sql"
)

type KeyValuesAPIRepository struct {
	database *sql.DB
}

func NewKeyValuesAPIRepository(database *sql.DB) KeyValuesAPIRepository {
	return KeyValuesAPIRepository{database}
}

func (r *KeyValuesAPIRepository) GetAllKeyValues() ([]models.APIKeyValue, error) {
	rows, err := r.database.Query(QueryAPIGetAllKeyValues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keyValues := []models.APIKeyValue{}
	for rows.Next() {
		var kv models.APIKeyValue
		err := rows.Scan(&kv.Key, &kv.Value)
		if err != nil {
			return nil, err
		}
		keyValues = append(keyValues, kv)
	}

	return keyValues, nil
}

func (r *KeyValuesAPIRepository) GetKeyValueByKey(key string) ([]models.APIKeyValue, error) {
	row := r.database.QueryRow(QueryAPIGetKeyValueByKey, key)

	var kv models.APIKeyValue
	err := row.Scan(&kv.Key, &kv.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty array if not found
			return []models.APIKeyValue{}, nil
		}
		return nil, err
	}

	return []models.APIKeyValue{kv}, nil
}

func (r *KeyValuesAPIRepository) GetKeyValuesStartingWith(prefix string) ([]models.APIKeyValue, error) {
	rows, err := r.database.Query(QueryAPIGetKeyValuesStartingWith, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keyValues := []models.APIKeyValue{}
	for rows.Next() {
		var kv models.APIKeyValue
		err := rows.Scan(&kv.Key, &kv.Value)
		if err != nil {
			return nil, err
		}
		keyValues = append(keyValues, kv)
	}

	return keyValues, nil
}
