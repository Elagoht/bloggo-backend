package user

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/middleware"
	"database/sql"
	"log"

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
	bucket, err := bucket.NewFileSystemBucket("users/avatars")
	if err != nil {
		log.Fatalln("User module cannot created file storage")
	}
	repository := NewUserRepository(database)
	service := NewUserService(repository, bucket)
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
			router.Post("/", module.Handler.UserCreate)
			router.Get("/", module.Handler.GetUsers)
			router.Get("/{id}", module.Handler.GetUserById)
			router.Get("/me", module.Handler.GetSelf)
			router.Patch("/me/avatar", module.Handler.UpdateSelfAvatar)
		})
}
