package auth

import (
	"bloggo/internal/module/auth/models"
	"database/sql"
)

type AuthRepository struct {
	database *sql.DB
}

func NewAuthRepository(database *sql.DB) AuthRepository {
	return AuthRepository{
		database,
	}
}

func (repository *AuthRepository) GetUserLoginDataByEmail(
	email string,
) (*models.UserLoginDetails, error) {
	row := repository.database.QueryRow(QueryUserLoginDataByEmail, email)

	var result = models.UserLoginDetails{}
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

func (repository *AuthRepository) GetUserLoginDataById(
	userId int64,
) (*models.UserLoginDetails, error) {
	row := repository.database.QueryRow(QueryUserLoginDataById, userId)

	var result = models.UserLoginDetails{}
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

func (repository *AuthRepository) GetUserPermissionsById(
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
