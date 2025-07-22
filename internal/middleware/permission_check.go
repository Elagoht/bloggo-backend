package middleware

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"net/http"
)

// Checks if the userRole in context has the required permission
func RequirePermission(
	permissionStore permissions.Store,
	requiredPermission string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			// Retrieve userRole from context
			role := request.Context().Value("userRole")
			if role == nil {
				// If userRole is not set, return 401 Unauthorized
				handlers.WriteError(writer, apierrors.NewAPIError(
					"role cannot acquired",
					apierrors.ErrUnauthorized,
				), http.StatusUnauthorized)
			}

			// FMT print all permission store and values' types

			// Check if userRole is a string and has the required permission
			if !permissionStore.HasPermission(role.(int64), requiredPermission) {
				handlers.WriteError(writer, apierrors.NewAPIError(
					"Insufficent permission",
					apierrors.ErrForbidden,
				), http.StatusForbidden)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
