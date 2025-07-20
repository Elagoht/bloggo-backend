package auth

import (
	"database/sql"

	"github.com/go-chi/chi"
)

type AuthModule struct {
	Handler    AuthHandler
	Service    AuthService
	Repository AuthRepository
}

func NewModule(database *sql.DB) AuthModule {
	repository := NewAuthRepository(database)
	service := NewAuthService(repository)
	handler := NewAuthHandler(service)

	return AuthModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module AuthModule) RegisterModule(
	database *sql.DB,
	router *chi.Mux,
) {
	router.Route("/auth", func(router chi.Router) {
		router.Post("/login", module.Handler.Login)
		router.Post("/refresh", module.Handler.Refresh)
		router.Post("/logout", module.Handler.Logout)
	})
}
