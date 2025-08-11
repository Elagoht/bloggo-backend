package category

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"bloggo/internal/utils/permission"

	"github.com/go-chi/chi"
)

type CategoryModule struct {
	Handler    CategoryHandler
	Service    CategoryService
	Repository CategoryRepository
}

func NewModule() CategoryModule {
	database := db.Get()
	repository := NewCategoryRepository(database)
	service := NewCategoryService(repository)
	handler := NewCategoryHandler(service)

	return CategoryModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module CategoryModule) RegisterModule(router *chi.Mux) {
	config := config.Get()
	permissionStore := permissions.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/categories",
		func(router chi.Router) {
			authority := permission.NewChecker(permissionStore)

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
