package middleware

import (
	"bloggo/internal/config"
	"bloggo/internal/utils/apierrors"
	"net/http"
)

// TrustedFrontendMiddleware validates the x-trusted-frontend header
func TrustedFrontendMiddleware(config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			trustedFrontendKey := request.Header.Get("x-trusted-frontend")

			if trustedFrontendKey == "" {
				apierrors.MapErrors(apierrors.ErrUnauthorized, writer, apierrors.HTTPErrorMapping{
					apierrors.ErrUnauthorized: {
						Message: "Missing x-trusted-frontend header",
						Status:  http.StatusUnauthorized,
					},
				})
				return
			}

			if trustedFrontendKey != config.TrustedFrontendKey {
				apierrors.MapErrors(apierrors.ErrUnauthorized, writer, apierrors.HTTPErrorMapping{
					apierrors.ErrUnauthorized: {
						Message: "Invalid x-trusted-frontend header",
						Status:  http.StatusUnauthorized,
					},
				})
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
