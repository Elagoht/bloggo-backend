package permission

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"net/http"
)

type Checker struct {
	permissions permissions.Store
}

func NewChecker(
	permissions permissions.Store,
) *Checker {
	return &Checker{
		permissions,
	}
}

func (checker *Checker) Require(
	wanted string,
	handler http.HandlerFunc,
) http.HandlerFunc {
	middleWare := middleware.RequirePermission(checker.permissions, wanted)
	return func(writer http.ResponseWriter, request *http.Request) {
		middleWare(handler).ServeHTTP(writer, request)
	}
}
