package keyvalue

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type KeyValueModule struct {
	Handler    KeyValueHandler
	Service    KeyValueService
	Repository KeyValueRepository
}

func NewModule() KeyValueModule {
	database := db.Get()
	permissionStore := permissions.Get()

	repository := NewKeyValueRepository(database)
	service := NewKeyValueService(repository, permissionStore)
	handler := NewKeyValueHandler(service)

	return KeyValueModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module KeyValueModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/key-values",
		func(router chi.Router) {
			router.Get("/", module.Handler.GetAll)
			router.Post("/bulk", module.Handler.BulkUpsert)
		},
	)
}
