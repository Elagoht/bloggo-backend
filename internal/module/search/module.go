package search

import (
	"bloggo/internal/db"

	"github.com/go-chi/chi"
)

type SearchModule struct {
	Handler    SearchHandler
	Service    SearchService
	Repository SearchRepository
}

func NewModule() SearchModule {
	database := db.Get()
	repository := NewSearchRepository(database)
	service := NewSearchService(repository)
	handler := NewSearchHandler(service)

	return SearchModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module SearchModule) RegisterModule(router *chi.Mux) {
	router.Route("/search", func(router chi.Router) {
		router.Get("/", module.Handler.Search)
	})
}