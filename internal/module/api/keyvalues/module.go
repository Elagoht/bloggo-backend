package keyvalues

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type KeyValuesAPIModule struct {
	Handler    KeyValuesAPIHandler
	Service    KeyValuesAPIService
	Repository KeyValuesAPIRepository
}

func NewModule() KeyValuesAPIModule {
	database := db.Get()

	repository := NewKeyValuesAPIRepository(database)
	service := NewKeyValuesAPIService(repository)
	handler := NewKeyValuesAPIHandler(service)

	return KeyValuesAPIModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module KeyValuesAPIModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.Route("/api/key-values", func(r chi.Router) {
		// All endpoints require trusted frontend header
		r.Use(middleware.TrustedFrontendMiddleware(&config))

		r.Get("/", module.Handler.ListKeyValues)
	})
}
