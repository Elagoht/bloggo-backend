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
		if !strings.HasPrefix(request.URL.Path, "/storage/") {
			writer.Header().Add("Content-Type", "application/json")
		}
		next.ServeHTTP(writer, request)
	})
}
