package auth

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/tokens"
	"database/sql"

	"github.com/go-chi/chi"
)

type AuthModule struct {
	Handler    AuthHandler
	Service    AuthService
	Repository AuthRepository
}

func NewModule(
	database *sql.DB,
	config *config.Config,
) AuthModule {
	refreshStore := tokens.GetStore()

	repository := NewAuthRepository(database)
	service := NewAuthService(repository, config, refreshStore)
	handler := NewAuthHandler(service, config)

	return AuthModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module AuthModule) RegisterModule(router *chi.Mux) {
	router.Route("/auth", func(router chi.Router) {
		router.Post("/login", module.Handler.Login)
		router.Post("/refresh", module.Handler.Refresh)
		router.Post("/logout", module.Handler.Logout)
	})
}
