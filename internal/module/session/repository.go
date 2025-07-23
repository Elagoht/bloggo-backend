package session

import (
	"bloggo/internal/module/session/models"
	"database/sql"
)

type SessionRepository struct {
	database *sql.DB
}

func NewSessionRepository(database *sql.DB) SessionRepository {
	return SessionRepository{
		database,
	}
}

func (repository *SessionRepository) GetUserLoginDataByEmail(
	email string,
) (*models.SessionCreateDetails, error) {
	row := repository.database.QueryRow(QuerySessionCreateDataByEmail, email)

	var result = models.SessionCreateDetails{}
	err := row.Scan(
		&result.UserId,
		&result.UserName,
		&result.RoleId,
		&result.RoleName,
		&result.PassphraseHash,
	)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (repository *SessionRepository) GetUserLoginDataById(
	userId int64,
) (*models.SessionCreateDetails, error) {
	row := repository.database.QueryRow(QuerySessionCreateDataById, userId)

	var result = models.SessionCreateDetails{}
	err := row.Scan(
		&result.UserId,
		&result.UserName,
		&result.RoleId,
		&result.RoleName,
		&result.PassphraseHash,
	)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (repository *SessionRepository) GetUserPermissionsById(
	userId int64,
) ([]string, error) {
	rows, err := repository.database.Query(QueryUserPermissionsById, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
