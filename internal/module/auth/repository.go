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
