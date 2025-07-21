package category

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"bloggo/internal/utils/permission"
	"database/sql"

	"github.com/go-chi/chi"
)

type CategoryModule struct {
	Handler     CategoryHandler
	Service     CategoryService
	Repository  CategoryRepository
	Config      *config.Config
	Permissions permissions.Store
}

func NewModule(
	database *sql.DB,
	config *config.Config,
	permissions permissions.Store,
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
			authority := permission.NewChecker(module.Permissions)

			router.Get("/",
				authority.Require("category:manage", module.Handler.GetCategories),
			)

			router.Get("/{slug}", module.Handler.GetCategoryBySlug)

			router.Post("/",
				authority.Require("category:manage", module.Handler.CategoryCreate),
			)

			router.Patch("/{slug}",
				authority.Require("category:manage", module.Handler.CategoryUpdate),
			)

			router.Delete("/{slug}",
				authority.Require("category:manage", module.Handler.CategoryDelete),
			)
		},
	)
}
