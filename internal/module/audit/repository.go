package audit

import (
	"bloggo/internal/module/audit/models"
	"bloggo/internal/utils/pagination"
	"database/sql"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) AuditRepository {
	return AuditRepository{
		db: db,
	}
}

func (repo *AuditRepository) GetAuditLogs(pagination *pagination.PaginationOptions) ([]models.ResponseAuditLog, error) {
	offset := 0
	limit := 20

	if pagination.Page != nil && pagination.Take != nil {
		offset = (*pagination.Page - 1) * *pagination.Take
		limit = *pagination.Take
	}

	rows, err := repo.db.Query(QueryGetAuditLogs, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auditLogs []models.ResponseAuditLog
	for rows.Next() {
		var log models.ResponseAuditLog
		err := rows.Scan(
			&log.Id,
			&log.UserId,
			&log.Entity,
			&log.EntityId,
			&log.Action,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		auditLogs = append(auditLogs, log)
	}

	return auditLogs, nil
}

func (repo *AuditRepository) GetAuditLogsCount() (int, error) {
	var count int
	err := repo.db.QueryRow(QueryGetAuditLogsCount).Scan(&count)
	return count, err
}