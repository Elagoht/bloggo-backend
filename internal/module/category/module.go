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
		router.Get("/", module.Handler.List)
		router.Post("/", module.Handler.Create)
		router.Put("/{id}", module.Handler.Update)
		router.Delete("/{id}", module.Handler.Delete)
	})
}
