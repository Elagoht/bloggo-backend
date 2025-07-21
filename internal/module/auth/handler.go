package auth

import (
	"bloggo/internal/config"
	"bloggo/internal/module/auth/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service AuthService
	config  *config.Config
}

func NewAuthHandler(
	service AuthService,
	config *config.Config,
) AuthHandler {
	return AuthHandler{
		service,
		config,
	}
}

func (handler *AuthHandler) Login(
	writer http.ResponseWriter,
	request *http.Request,
) {
	body, ok := handlers.BindAndValidate[*models.RequestLogin](writer, request)
	if !ok {
		return
	}

	accessToken, refreshToken, err := handler.service.GenerateTokens(body)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	handler.sendTokens(writer, accessToken, refreshToken)
}

func (handler *AuthHandler) Refresh(
	writer http.ResponseWriter,
	request *http.Request,
) {
	// Get refresh token
	refreshCookie, err := request.Cookie("refreshToken")
	if err != nil {
		apierrors.MapErrors(apierrors.ErrUnauthorized, writer, nil)
		return
	}

	// Refresh all tokens
	accessToken, refreshToken, err := handler.service.RefreshTokens(
		refreshCookie.Value,
	)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	handler.sendTokens(writer, accessToken, refreshToken)
}

func (handler *AuthHandler) Logout(
	writer http.ResponseWriter,
	request *http.Request,
) {
	// Get refresh token
	refreshCookie, err := request.Cookie("refreshToken")
	if err != nil {
		apierrors.MapErrors(apierrors.ErrUnauthorized, writer, nil)
		return
	}

	// Revoke refresh token from store
	handler.service.RevokeRefreshToken(refreshCookie.Value)

	// Remove refresh token from client
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   -1, // Expire immediately
	}
	http.SetCookie(writer, &cookie)
	writer.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) sendTokens(
	writer http.ResponseWriter,
	accessToken string,
	refreshToken string,
) {
	// Set refresh token as an HTTP-Only cookie
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   handler.config.RefreshTokenDuration,
	}
	http.SetCookie(writer, &cookie)

	// Write access token to response body
	response := models.ResponseAccessToken{
		AccessToken: accessToken,
	}
	json.NewEncoder(writer).Encode(response)
}
