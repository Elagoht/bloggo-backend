package authors

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type AuthorsAPIModule struct {
	Handler    AuthorsAPIHandler
	Service    AuthorsAPIService
	Repository AuthorsAPIRepository
}

func NewModule() AuthorsAPIModule {
	database := db.Get()

	repository := NewAuthorsAPIRepository(database)
	service := NewAuthorsAPIService(repository)
	handler := NewAuthorsAPIHandler(service)

	return AuthorsAPIModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module AuthorsAPIModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.Route("/api/authors", func(r chi.Router) {
		// All endpoints require trusted frontend header
		r.Use(middleware.TrustedFrontendMiddleware(&config))

		r.Get("/", module.Handler.ListAuthors)
		r.Get("/{id}", module.Handler.GetAuthorById)
	})
}
