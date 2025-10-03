package webhook

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type WebhookModule struct {
	Handler    WebhookHandler
	Service    WebhookService
	Repository WebhookRepository
}

func NewModule() WebhookModule {
	database := db.Get()
	permissionStore := permissions.Get()

	repository := NewWebhookRepository(database)
	service := NewWebhookService(repository, permissionStore)
	handler := NewWebhookHandler(service)

	return WebhookModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module WebhookModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/webhook",
		func(router chi.Router) {
			router.Get("/config", module.Handler.GetConfig)
			router.Put("/config", module.Handler.UpdateConfig)

			router.Get("/headers", module.Handler.GetHeaders)
			router.Put("/headers", module.Handler.BulkUpsertHeaders)

			router.Post("/fire", module.Handler.ManualFire)

			router.Get("/requests", module.Handler.GetRequests)
			router.Get("/requests/{id}", module.Handler.GetRequestByID)
		},
	)
}
