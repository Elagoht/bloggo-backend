package middleware

import (
	permissionstore "bloggo/internal/infrastructure/permission_store"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"net/http"
)

// Checks if the userRole in context has the required permission
func RequirePermission(
	permissionStore permissionstore.PermissionStore,
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

			// Check if userRole is a string and has the required permission
			roleStr, ok := role.(string)
			if !ok || !permissionStore.HasPermission(roleStr, requiredPermission) {
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
