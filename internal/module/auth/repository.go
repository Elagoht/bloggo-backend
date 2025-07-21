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
		&result.RoleId,
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
	row := repository.database.QueryRow("SELECT id, role_id, passphrase_hash FROM users WHERE id = ? AND deleted_at IS NULL;", userId)

	var result = models.UserLoginDetails{}
	err := row.Scan(
		&result.UserId,
		&result.RoleId,
		&result.PassphraseHash,
	)
	if err != nil {
		return nil, err
	}

	return &result, err
}
