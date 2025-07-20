package category

import (
	"database/sql"

	"github.com/go-chi/chi"
)

type CategoryModule struct {
	Handler    CategoryHandler
	Service    CategoryService
	Repository CategoryRepository
}

func NewModule(database *sql.DB) CategoryModule {
	repository := NewCategoryRepository(database)
	service := NewCategoryService(repository)
	handler := NewCategoryHandler(service)

	return CategoryModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module CategoryModule) RegisterModule(
	database *sql.DB,
	router *chi.Mux,
) {
	router.Route("/categories", func(router chi.Router) {
		router.Get("/", module.Handler.GetCategories)
		router.Get("/{slug}", module.Handler.GetCategoryBySlug)
		router.Post("/", module.Handler.CategoryCreate)
		router.Patch("/{slug}", module.Handler.CategoryUpdate)
		router.Delete("/{slug}", module.Handler.CategoryDelete)
	})
}
