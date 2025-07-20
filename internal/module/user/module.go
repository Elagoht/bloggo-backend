package user

import (
	"database/sql"

	"github.com/go-chi/chi"
)

type UserModule struct {
	Handler    UserHandler
	Service    UserService
	Repository UserRepository
}

func NewModule(database *sql.DB) UserModule {
	repository := NewUserRepository(database)
	service := NewUserService(repository)
	handler := NewUserHandler(service)

	return UserModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module UserModule) RegisterModule(
	database *sql.DB,
	router *chi.Mux,
) {
	router.Route("/users", func(router chi.Router) {
		router.Get("/", module.Handler.GetUsers)
		router.Get("/{id}", module.Handler.GetUserById)
		// router.Post("/", module.Handler.UserCreate)
		// router.Patch("/{id}", module.Handler.UserUpdate)
		// router.Delete("/{id}", module.Handler.UserDelete)
	})
}
