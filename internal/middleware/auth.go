package middleware

import (
	"bloggo/internal/config"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Reads the JWT access token from the Authorization header, validates it, and sets userRole in the context.
func AuthMiddleware(
	configuration *config.Config,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			// Extract the Authorization header
			header := request.Header.Get("Authorization")

			// Check if the header is present and starts with "Bearer "
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				handlers.WriteError(
					writer,
					apierrors.NewAPIError(
						"Missing or invalid Authorization header",
						apierrors.ErrUnauthorized,
					),
					http.StatusUnauthorized,
				)
				return
			}

			// Remove "Bearer " prefix to get the token string
			tokenString := strings.TrimPrefix(header, "Bearer ")
			// Decode the JWT secret from base64
			key, err := base64.RawURLEncoding.DecodeString(configuration.JWTSecret)
			if err != nil {
				// If the secret can't be decoded, return 500 Internal Server Error
				handlers.WriteError(
					writer,
					apierrors.NewAPIError(
						"Server misconfiguration",
						err,
					),
					http.StatusInternalServerError,
				)
				return
			}

			// Parse and validate the JWT token
			claims := jwt.MapClaims{}
			_, err = jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (any, error) { return key, nil },
			)
			if err != nil {
				// If the token is invalid or expired, return 401 Unauthorized
				handlers.WriteError(
					writer,
					apierrors.NewAPIError(
						"Invalid or expired token",
						apierrors.ErrUnauthorized,
					),
					http.StatusUnauthorized,
				)
				return
			}

			// Extract the role id (rid) and user id (uid) from the claims (JWT numbers are float64)
			rid, ok := claims["rid"].(float64)
			if !ok {
				// If the role is not found, return 401 Unauthorized
				handlers.WriteError(
					writer,
					apierrors.NewAPIError(
						"Role not found in token",
						apierrors.ErrUnauthorized,
					),
					http.StatusUnauthorized,
				)
				return
			}
			uid, ok := claims["uid"].(float64)
			if !ok {
				// If the user id is not found, return 401 Unauthorized
				handlers.WriteError(
					writer,
					apierrors.NewAPIError(
						"User id not found in token",
						apierrors.ErrUnauthorized,
					),
					http.StatusUnauthorized,
				)
				return
			}

			// Set userRole in the request context as int64
			newContext := context.WithValue(
				request.Context(),
				handlers.TokenRoleId,
				int64(rid),
			)
			newContext = context.WithValue(newContext, handlers.TokenUserId, int64(uid))
			// Call the next handler with the updated context

			next.ServeHTTP(writer, request.WithContext(newContext))
		})
	}
}
