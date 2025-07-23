package user

import (
	"bloggo/internal/config"
	"bloggo/internal/middleware"
	"database/sql"

	"github.com/go-chi/chi"
)

type UserModule struct {
	Handler    UserHandler
	Service    UserService
	Repository UserRepository
	Config     *config.Config
}

func NewModule(
	database *sql.DB,
	config *config.Config,
) UserModule {
	repository := NewUserRepository(database)
	service := NewUserService(repository)
	handler := NewUserHandler(service)

	return UserModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
		Config:     config,
	}
}

func (module UserModule) RegisterModule(router *chi.Mux) {
	router.With(middleware.AuthMiddleware(module.Config)).Route(
		"/users",
		func(router chi.Router) {
			router.Get("/me", module.Handler.GetSelf)
			router.Get("/", module.Handler.GetUsers)
			router.Get("/{id}", module.Handler.GetUserById)
			router.Post("/", module.Handler.UserCreate)
		})
}
