package keyvalue

import (
	"bloggo/internal/module/keyvalue/models"
	"database/sql"
)

type KeyValueRepository struct {
	database *sql.DB
}

func NewKeyValueRepository(database *sql.DB) KeyValueRepository {
	return KeyValueRepository{
		database,
	}
}

func (repository *KeyValueRepository) GetAll() ([]models.KeyValue, error) {
	rows, err := repository.database.Query(QueryGetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.KeyValue{}
	for rows.Next() {
		var item models.KeyValue
		err := rows.Scan(
			&item.Key,
			&item.Value,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repository *KeyValueRepository) Upsert(key, value string) error {
	statement, err := repository.database.Prepare(QueryUpsert)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(key, value)
	return err
}

func (repository *KeyValueRepository) BulkUpsert(items []models.RequestKeyValueUpsert) error {
	tx, err := repository.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statement, err := tx.Prepare(QueryUpsert)
	if err != nil {
		return err
	}
	defer statement.Close()

	for _, item := range items {
		_, err := statement.Exec(item.Key, item.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repository *KeyValueRepository) DeleteAll() error {
	_, err := repository.database.Exec(QueryDeleteAll)
	return err
}
