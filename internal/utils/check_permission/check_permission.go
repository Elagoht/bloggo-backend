package checkpermission

import (
	permissionstore "bloggo/internal/infrastructure/permission_store"
	"bloggo/internal/middleware"
	"net/http"
)

type Checker struct {
	permissions permissionstore.PermissionStore
}

func NewChecker(
	permissions permissionstore.PermissionStore,
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
