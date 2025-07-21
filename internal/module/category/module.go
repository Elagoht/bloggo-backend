package category

import (
	"bloggo/internal/config"
	permissionstore "bloggo/internal/infrastructure/permission_store"
	"bloggo/internal/middleware"
	checkpermission "bloggo/internal/utils/check_permission"
	"database/sql"

	"github.com/go-chi/chi"
)

type CategoryModule struct {
	Handler     CategoryHandler
	Service     CategoryService
	Repository  CategoryRepository
	Config      *config.Config
	Permissions permissionstore.PermissionStore
}

func NewModule(
	database *sql.DB,
	config *config.Config,
	permissions permissionstore.PermissionStore,
) CategoryModule {
	repository := NewCategoryRepository(database)
	service := NewCategoryService(repository)
	handler := NewCategoryHandler(service)

	return CategoryModule{
		Handler:     handler,
		Service:     service,
		Repository:  repository,
		Config:      config,
		Permissions: permissions,
	}
}

func (module CategoryModule) RegisterModule(router *chi.Mux) {
	router.With(middleware.AuthMiddleware(module.Config)).Route(
		"/categories",
		func(router chi.Router) {
			permission := checkpermission.NewChecker(module.Permissions)

			router.Get("/",
				permission.Require("category:manage", module.Handler.GetCategories),
			)

			router.Get("/{slug}", module.Handler.GetCategoryBySlug)

			router.Post("/",
				permission.Require("category:manage", module.Handler.CategoryCreate),
			)

			router.Patch("/{slug}",
				permission.Require("category:manage", module.Handler.CategoryUpdate),
			)

			router.Delete("/{slug}",
				permission.Require("category:manage", module.Handler.CategoryDelete),
			)
		},
	)
}
