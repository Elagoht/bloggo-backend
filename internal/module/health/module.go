package health

import (
	"github.com/go-chi/chi"
)

type HealthModule struct {
	Handler HealthHandler
}

func NewModule() HealthModule {
	handler := NewHealthHandler()

	return HealthModule{
		Handler: handler,
	}
}

func (module HealthModule) RegisterModule(router *chi.Mux) {
	router.Route(
		"/health",
		func(router chi.Router) {
			router.Get("/", module.Handler.CheckHealth)
		},
	)
}
