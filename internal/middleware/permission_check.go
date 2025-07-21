package middleware

import (
	permissionstore "bloggo/internal/infrastructure/permission_store"
	"net/http"
)

func RequirePermission(
	permissionStore permissionstore.PermissionStore,
	requiredPermission string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Kullanıcının rolünü context'ten veya header'dan al
			role := r.Context().Value("userRole")
			if role == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			roleStr, ok := role.(string)
			if !ok || !permissionStore.HasPermission(roleStr, requiredPermission) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Yetkisi varsa devam et
			next.ServeHTTP(w, r)
		})
	}
}
