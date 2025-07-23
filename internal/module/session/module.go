package session

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/tokens"
	"database/sql"

	"github.com/go-chi/chi"
)

type SessionModule struct {
	Handler    SessionHandler
	Service    SessionService
	Repository SessionRepository
}

func NewModule(
	database *sql.DB,
	config *config.Config,
) SessionModule {
	refreshStore := tokens.GetStore()

	repository := NewSessionRepository(database)
	service := NewSessionService(repository, config, refreshStore)
	handler := NewSessionHandler(service, config)

	return SessionModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module SessionModule) RegisterModule(router *chi.Mux) {
	router.Route("/session", func(router chi.Router) {
		router.Post("/", module.Handler.CreateSession)
		router.Post("/refresh", module.Handler.RefreshSession)
		router.Delete("/", module.Handler.DeleteSession)
	})
}
