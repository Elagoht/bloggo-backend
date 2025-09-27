package session

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/tokens"
	"bloggo/internal/module/audit"
	auditmodels "bloggo/internal/module/audit/models"
	"bloggo/internal/module/session/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/cryptography"
)

type UserService interface {
	UpdateLastLogin(userId int64) error
}

type SessionService struct {
	repository   SessionRepository
	config       *config.Config
	refreshStore tokens.Store
	userService  UserService
}

func NewSessionService(
	repository SessionRepository,
	config *config.Config,
	refreshStore tokens.Store,
	userService UserService,
) SessionService {
	return SessionService{
		repository,
		config,
		refreshStore,
		userService,
	}
}

func (service *SessionService) CreateSession(
	model *models.RequestSessionCreate,
) (session *models.ResponseSession, refreshToken string, err error) {
	// Compare passphrase hashes
	details, err := service.repository.GetUserLoginDataByEmail(model.Email)
	if err != nil {
		// Not sending "resource not found" error
		// Do not allow hackers to brute force to
		// find registered emails
		return nil, "", apierrors.ErrUnauthorized
	}

	if !cryptography.ComparePassphrase(
		details.PassphraseHash,
		model.Passphrase,
	) {
		return nil, "", apierrors.ErrUnauthorized
	}

	// Generate tokens
	accessToken, err := cryptography.GenerateJWT(
		model.Email,
		details.UserId,
		details.RoleId,
		service.config.JWTSecret,
		service.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, "", err
	}

	refreshToken = cryptography.GenerateUniqueId()

	// Get user permissions
	permissions, err := service.repository.GetUserPermissionsById(details.UserId)
	if err != nil {
		return nil, "", err
	}

	// Set refresh token to Refresh Token Store
	// to be able to revoke sessions by hand
	service.refreshStore.Set(
		refreshToken,
		details.UserId,
		service.config.RefreshTokenDuration,
	)

	// Update last login timestamp
	err = service.userService.UpdateLastLogin(details.UserId)
	if err != nil {
		// Log the error but don't fail the login process
		// This is not critical for authentication to succeed
	}

	sessionData := &models.ResponseSession{
		AccessToken: accessToken,
		Name:        details.UserName,
		Role:        details.RoleName,
		Permissions: permissions,
	}

	// Log login action
	audit.LogAuthAction(&details.UserId, auditmodels.ActionLogin, nil, nil)

	return sessionData, refreshToken, nil
}

func (service *SessionService) RefreshSession(
	refreshToken string,
) (session *models.ResponseSession, rotatedRefreshToken string, err error) {
	userId, found := service.refreshStore.Get(refreshToken)
	if !found {
		return nil, "", apierrors.ErrUnauthorized
	}

	details, err := service.repository.GetUserLoginDataById(userId)
	if err != nil {
		return nil, "", apierrors.ErrUnauthorized
	}

	// Generate new access token
	accessToken, err := cryptography.GenerateJWT(
		"", // Email is not available here, can be added if needed
		details.UserId,
		details.RoleId,
		service.config.JWTSecret,
		service.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, "", err
	}

	// Get user permissions
	permissions, err := service.repository.GetUserPermissionsById(details.UserId)
	if err != nil {
		return nil, "", err
	}

	// Rotate refresh token
	newRefreshToken := cryptography.GenerateUniqueId()

	// Set new refresh token to Refresh Token Store
	service.refreshStore.Set(
		newRefreshToken,
		details.UserId,
		service.config.RefreshTokenDuration,
	)
	// Revoke old refresh token
	service.refreshStore.Delete(refreshToken)

	sessionData := &models.ResponseSession{
		AccessToken: accessToken,
		Name:        details.UserName,
		Role:        details.RoleName,
		Permissions: permissions,
	}

	return sessionData, newRefreshToken, nil
}

func (service *SessionService) RevokeSession(
	refreshToken string,
) {
	// Get user ID before deleting token for audit logging
	userId, found := service.refreshStore.Get(refreshToken)

	// Revoke refresh token
	service.refreshStore.Delete(refreshToken)

	// Log logout action if we found the user
	if found {
		audit.LogAuthAction(&userId, auditmodels.ActionLogout, nil, nil)
	}
}
