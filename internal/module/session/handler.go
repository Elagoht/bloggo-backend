package session

import (
	"bloggo/internal/config"
	"bloggo/internal/module/session/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type SessionHandler struct {
	service SessionService
	config  *config.Config
}

func NewSessionHandler(
	service SessionService,
	config *config.Config,
) SessionHandler {
	return SessionHandler{
		service,
		config,
	}
}

func (handler *SessionHandler) CreateSession(
	writer http.ResponseWriter,
	request *http.Request,
) {
	body, ok := handlers.BindAndValidate[*models.RequestSessionCreate](writer, request)
	if !ok {
		return
	}

	session, refreshToken, err := handler.service.CreateSession(body)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	handler.sendSession(writer, session, refreshToken)
}

func (handler *SessionHandler) RefreshSession(
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
	session, newRefreshToken, err := handler.service.RefreshSession(
		refreshCookie.Value,
	)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	handler.sendSession(writer, session, newRefreshToken)
}

func (handler *SessionHandler) DeleteSession(
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
	handler.service.RevokeSession(refreshCookie.Value)

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

func (handler *SessionHandler) sendSession(
	writer http.ResponseWriter,
	session *models.ResponseSession,
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
	response := models.ResponseSession{
		AccessToken: session.AccessToken,
		Name:        session.Name,
		Role:        session.Role,
		Permissions: session.Permissions,
	}
	json.NewEncoder(writer).Encode(response)
}
