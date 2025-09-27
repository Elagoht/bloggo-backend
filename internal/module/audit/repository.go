package audit

import (
	"bloggo/internal/module/audit/models"
	"database/sql"
	"encoding/json"
	"strings"
)

type AuditRepository struct {
	database *sql.DB
}

func NewAuditRepository(database *sql.DB) AuditRepository {
	return AuditRepository{
		database,
	}
}

func (repository *AuditRepository) LogAction(entry *models.AuditLogEntry) error {
	var metadataJSON *string

	// Convert metadata to JSON string
	if entry.Metadata != nil {
		if jsonBytes, err := json.Marshal(entry.Metadata); err == nil {
			jsonStr := string(jsonBytes)
			metadataJSON = &jsonStr
		}
	}

	_, err := repository.database.Exec(
		QueryInsertAuditLog,
		entry.UserID,
		entry.EntityType,
		entry.EntityID,
		entry.Action,
		metadataJSON,
	)

	return err
}

func (repository *AuditRepository) GetAuditLogs(limit, offset int) ([]models.AuditLogResponse, error) {
	rows, err := repository.database.Query(QueryGetAuditLogs, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanAuditLogs(rows)
}

func (repository *AuditRepository) GetAuditLogsByEntity(entityType string, entityID int64, limit, offset int) ([]models.AuditLogResponse, error) {
	rows, err := repository.database.Query(QueryGetAuditLogsByEntity, entityType, entityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanAuditLogs(rows)
}

func (repository *AuditRepository) GetAuditLogsByUser(userID int64, limit, offset int) ([]models.AuditLogResponse, error) {
	rows, err := repository.database.Query(QueryGetAuditLogsByUser, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanAuditLogs(rows)
}

func (repository *AuditRepository) CountAuditLogs() (int, error) {
	row := repository.database.QueryRow(QueryCountAuditLogs)
	var count int
	err := row.Scan(&count)
	return count, err
}

func (repository *AuditRepository) CountAuditLogsByEntity(entityType string, entityID int64) (int, error) {
	row := repository.database.QueryRow(QueryCountAuditLogsByEntity, entityType, entityID)
	var count int
	err := row.Scan(&count)
	return count, err
}

func (repository *AuditRepository) CountAuditLogsByUser(userID int64) (int, error) {
	row := repository.database.QueryRow(QueryCountAuditLogsByUser, userID)
	var count int
	err := row.Scan(&count)
	return count, err
}

func (repository *AuditRepository) GetAuditLogsWithFilters(limit, offset int, userIDs []int64, entityTypes, actions []string) ([]models.AuditLogResponse, error) {
	var query string = `
	SELECT
		al.id, al.user_id, u.name as user_name,
		al.entity_type, al.entity_id, al.action,
		al.metadata,
		al.created_at,
		CASE
			WHEN al.entity_type = 'user' THEN (SELECT name FROM users WHERE id = al.entity_id)
			WHEN al.entity_type = 'category' THEN (SELECT name FROM categories WHERE id = al.entity_id)
			WHEN al.entity_type = 'tag' THEN (SELECT name FROM tags WHERE id = al.entity_id)
			WHEN al.entity_type = 'post' THEN (SELECT pv.title FROM post_versions pv JOIN posts p ON p.current_version_id = pv.id WHERE p.id = al.entity_id)
			WHEN al.entity_type = 'post_version' THEN (SELECT title FROM post_versions WHERE id = al.entity_id)
			ELSE NULL
		END as entity_name
	FROM audit_logs al
	LEFT JOIN users u ON al.user_id = u.id
	WHERE 1=1`

	var args []interface{}

	// Handle user ID array filter (OR logic - show logs from any of these users)
	if len(userIDs) > 0 {
		placeholders := make([]string, len(userIDs))
		for i, userID := range userIDs {
			placeholders[i] = "?"
			args = append(args, userID)
		}
		query += " AND al.user_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	// Handle entity type array filter (OR logic - show logs for any of these entity types)
	if len(entityTypes) > 0 {
		placeholders := make([]string, len(entityTypes))
		for i, entityType := range entityTypes {
			placeholders[i] = "?"
			args = append(args, entityType)
		}
		query += " AND al.entity_type IN (" + strings.Join(placeholders, ",") + ")"
	}

	// Handle action array filter (OR logic - show logs with any of these actions)
	if len(actions) > 0 {
		placeholders := make([]string, len(actions))
		for i, action := range actions {
			placeholders[i] = "?"
			args = append(args, action)
		}
		query += " AND al.action IN (" + strings.Join(placeholders, ",") + ")"
	}


	query += " ORDER BY al.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := repository.database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanAuditLogs(rows)
}

func (repository *AuditRepository) CountAuditLogsWithFilters(userIDs []int64, entityTypes, actions []string) (int, error) {
	var query string = "SELECT COUNT(*) FROM audit_logs WHERE 1=1"
	var args []interface{}

	// Handle user ID array filter
	if len(userIDs) > 0 {
		placeholders := make([]string, len(userIDs))
		for i, userID := range userIDs {
			placeholders[i] = "?"
			args = append(args, userID)
		}
		query += " AND user_id IN (" + strings.Join(placeholders, ",") + ")"
	}

	// Handle entity type array filter
	if len(entityTypes) > 0 {
		placeholders := make([]string, len(entityTypes))
		for i, entityType := range entityTypes {
			placeholders[i] = "?"
			args = append(args, entityType)
		}
		query += " AND entity_type IN (" + strings.Join(placeholders, ",") + ")"
	}

	// Handle action array filter
	if len(actions) > 0 {
		placeholders := make([]string, len(actions))
		for i, action := range actions {
			placeholders[i] = "?"
			args = append(args, action)
		}
		query += " AND action IN (" + strings.Join(placeholders, ",") + ")"
	}


	row := repository.database.QueryRow(query, args...)
	var count int
	err := row.Scan(&count)
	return count, err
}

func (repository *AuditRepository) scanAuditLogs(rows *sql.Rows) ([]models.AuditLogResponse, error) {
	var logs []models.AuditLogResponse

	for rows.Next() {
		var log models.AuditLogResponse
		var metadataJSON sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.UserName,
			&log.EntityType,
			&log.EntityID,
			&log.Action,
			&metadataJSON,
			&log.CreatedAt,
			&log.EntityName,
		)

		if err != nil {
			return nil, err
		}

		// Parse JSON strings back to maps
		if metadataJSON.Valid {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(metadataJSON.String), &metadata); err == nil {
				log.Metadata = &metadata
			}
		}

		logs = append(logs, log)
	}

	return logs, rows.Err()
}