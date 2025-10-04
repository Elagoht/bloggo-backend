package webhook

import (
	"bloggo/internal/module/webhook/models"
	"database/sql"
)

type WebhookRepository struct {
	database *sql.DB
}

func NewWebhookRepository(database *sql.DB) WebhookRepository {
	return WebhookRepository{
		database,
	}
}

// Config methods
func (repository *WebhookRepository) GetConfig() (*models.WebhookConfig, error) {
	var config models.WebhookConfig
	err := repository.database.QueryRow(QueryGetConfig).Scan(
		&config.ID,
		&config.URL,
		&config.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (repository *WebhookRepository) UpsertConfig(url string) error {
	statement, err := repository.database.Prepare(QueryUpsertConfig)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(url)
	return err
}

// Header methods
func (repository *WebhookRepository) GetAllHeaders() ([]models.WebhookHeader, error) {
	rows, err := repository.database.Query(QueryGetAllHeaders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	headers := []models.WebhookHeader{}
	for rows.Next() {
		var header models.WebhookHeader
		err := rows.Scan(
			&header.ID,
			&header.Key,
			&header.Value,
			&header.CreatedAt,
			&header.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		headers = append(headers, header)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return headers, nil
}

func (repository *WebhookRepository) BulkUpsertHeaders(items []models.RequestHeaderUpsert) error {
	tx, err := repository.database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete all existing headers first
	_, err = tx.Exec(QueryDeleteAllHeaders)
	if err != nil {
		return err
	}

	// Insert new headers
	statement, err := tx.Prepare(QueryUpsertHeader)
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

// Request methods
func (repository *WebhookRepository) InsertRequest(req *models.WebhookRequest) (int64, error) {
	statement, err := repository.database.Prepare(QueryInsertRequest)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(
		req.Event,
		req.Entity,
		req.EntityID,
		req.Slug,
		req.RequestBody,
		req.ResponseStatus,
		req.ResponseBody,
		req.AttemptCount,
		req.ErrorMessage,
		req.WebhookURL,
		req.WebhookHeaders,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (repository *WebhookRepository) UpdateRequest(id int64, responseStatus *int, responseBody *string, attemptCount int, errorMessage *string) error {
	statement, err := repository.database.Prepare(QueryUpdateRequest)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(responseStatus, responseBody, attemptCount, errorMessage, id)
	return err
}

func (repository *WebhookRepository) GetAllRequests(limit, offset int) ([]models.WebhookRequest, error) {
	rows, err := repository.database.Query(QueryGetAllRequests, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanRequests(rows)
}

func (repository *WebhookRepository) GetRequestsBySearch(search string, limit, offset int) ([]models.WebhookRequest, error) {
	searchPattern := "%" + search + "%"
	rows, err := repository.database.Query(QueryGetRequestsBySearch, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanRequests(rows)
}

func (repository *WebhookRepository) GetRequestByID(id int) (*models.WebhookRequest, error) {
	var req models.WebhookRequest
	err := repository.database.QueryRow(QueryGetRequestByID, id).Scan(
		&req.ID,
		&req.Event,
		&req.Entity,
		&req.EntityID,
		&req.Slug,
		&req.RequestBody,
		&req.ResponseStatus,
		&req.ResponseBody,
		&req.AttemptCount,
		&req.ErrorMessage,
		&req.WebhookURL,
		&req.WebhookHeaders,
		&req.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (repository *WebhookRepository) CountRequests() (int, error) {
	var count int
	err := repository.database.QueryRow(QueryCountRequests).Scan(&count)
	return count, err
}

func (repository *WebhookRepository) CountRequestsBySearch(search string) (int, error) {
	searchPattern := "%" + search + "%"
	var count int
	err := repository.database.QueryRow(QueryCountRequestsBySearch, searchPattern, searchPattern).Scan(&count)
	return count, err
}

// Helper to scan requests
func (repository *WebhookRepository) scanRequests(rows *sql.Rows) ([]models.WebhookRequest, error) {
	requests := []models.WebhookRequest{}
	for rows.Next() {
		var req models.WebhookRequest
		err := rows.Scan(
			&req.ID,
			&req.Event,
			&req.Entity,
			&req.EntityID,
			&req.Slug,
			&req.RequestBody,
			&req.ResponseStatus,
			&req.ResponseBody,
			&req.AttemptCount,
			&req.ErrorMessage,
			&req.WebhookURL,
			&req.WebhookHeaders,
			&req.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}
