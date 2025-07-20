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

	accessToken, refreshToken, err := handler.service.Login(body)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

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

func (handler *AuthHandler) Refresh(
	writer http.ResponseWriter,
	request *http.Request,
) {
	// TODO: Implement refresh token logic
	// 1. Read refresh token from cookie
	// 2. Validate refresh token
	// 3. Issue new access token (and possibly new refresh token)
	// 4. Return new access token in response
}

func (handler *AuthHandler) Logout(
	writer http.ResponseWriter,
	request *http.Request,
) {
	// TODO: Implement logout logic
	// 1. Invalidate refresh token (if using server-side storage)
	// 2. Remove refresh token cookie
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
