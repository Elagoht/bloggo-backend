package auth

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/tokens"
	"bloggo/internal/module/auth/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/cryptography"
)

type AuthService struct {
	repository   AuthRepository
	config       *config.Config
	refreshStore tokens.Store
}

func NewAuthService(
	repository AuthRepository,
	config *config.Config,
	refreshStore tokens.Store,

) AuthService {
	return AuthService{
		repository,
		config,
		refreshStore,
	}
}

func (service *AuthService) GenerateTokens(
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

	// Set refresh token to Refresh Token Store
	// to be able to revoke sessions by hand
	service.refreshStore.Set(
		refreshToken,
		details.UserId,
		service.config.RefreshTokenDuration,
	)

	return accessToken, refreshToken, nil
}

func (service *AuthService) RefreshTokens(
	refreshToken string,
) (accessToken string, rotatedRefreshToken string, err error) {
	userId, found := service.refreshStore.Get(refreshToken)
	if !found {
		return "", "", apierrors.ErrUnauthorized
	}

	details, err := service.repository.GetUserLoginDataById(userId)
	if err != nil {
		return "", "", apierrors.ErrUnauthorized
	}

	// Generate new access token
	accessToken, err = cryptography.GenerateJWT(
		"", // Email is not available here, can be added if needed
		details.UserId,
		details.RoleId,
		service.config.JWTSecret,
		service.config.AccessTokenDuration,
	)
	if err != nil {
		return "", "", err
	}

	// Rotate refresh token
	newRefreshToken, err := cryptography.GenerateUniqueId()
	if err != nil {
		return "", "", err
	}

	// Set new refresh token to Refresh Token Store
	service.refreshStore.Set(
		newRefreshToken,
		details.UserId,
		service.config.RefreshTokenDuration,
	)
	// Revoke old refresh token
	service.refreshStore.Delete(refreshToken)

	return accessToken, newRefreshToken, nil
}

func (service *AuthService) RevokeRefreshToken(
	refreshToken string,
) {
	// Revoke refresh token
	service.refreshStore.Delete(refreshToken)
}
