package middleware

import (
	"net/http"
	"strings"
)

func ResponseJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		// Skip JSON content type for static files and file storage
		path := request.URL.Path
		isStaticRoute := strings.HasPrefix(path, "/bucket/") ||
			strings.HasPrefix(path, "/assets/") ||
			path == "/" ||
			(!strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/internal/"))

		if !isStaticRoute {
			writer.Header().Add("Content-Type", "application/json")
		}
		next.ServeHTTP(writer, request)
	})
}
