package auth

import (
	"bloggo/internal/config"
	"bloggo/internal/module/auth/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/cryptography"
)

type AuthService struct {
	repository AuthRepository
	config     *config.Config
}

func NewAuthService(
	repository AuthRepository,
	config *config.Config,
) AuthService {
	return AuthService{
		repository,
		config,
	}
}

func (service *AuthService) Login(
	model *models.RequestLogin,
) (accessToken string, refreshToken string, err error) {
	// Compare passphrase hashes
	details, err := service.repository.GetUserLoginDataByEmail(model.Email)
	if err != nil {
		return "", "", err
	}

	if !cryptography.ComparePassphrase(
		details.PassphraseHash,
		model.Passphrase,
	) {
		return "", "", apierrors.ErrUnauthorized
	}

	// Generate tokens
	accessToken, err = cryptography.GenerateJWT(
		model.Email,
		details.UserId,
		details.RoleId,
		service.config.JWTSecret,
		service.config.AccessTokenDuration,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = cryptography.GenerateUniqueId()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
