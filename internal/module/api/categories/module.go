package categories

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type CategoriesAPIModule struct {
	Handler    CategoriesAPIHandler
	Service    CategoriesAPIService
	Repository CategoriesAPIRepository
}

func NewModule() CategoriesAPIModule {
	database := db.Get()

	repository := NewCategoriesAPIRepository(database)
	service := NewCategoriesAPIService(repository)
	handler := NewCategoriesAPIHandler(service)

	return CategoriesAPIModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module CategoriesAPIModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.Route("/api/categories", func(r chi.Router) {
		// All endpoints require trusted frontend header
		r.Use(middleware.TrustedFrontendMiddleware(&config))

		r.Get("/", module.Handler.ListCategories)
		r.Get("/{slug}", module.Handler.GetCategoryBySlug)
	})
}
