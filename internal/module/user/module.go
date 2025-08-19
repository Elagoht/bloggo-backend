package user

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/middleware"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"log"

	"github.com/go-chi/chi"
)

type UserModule struct {
	Handler    UserHandler
	Service    UserService
	Repository UserRepository
}

func NewModule() UserModule {
	database := db.Get()
	bucket, err := bucket.NewFileSystemBucket("users/avatars")
	if err != nil {
		log.Fatalln("User module cannot created file storage")
	}
	imageValidator := validatefile.NewImageValidator(5 << 20) // 5MB
	avatarResizer := transformfile.NewImageTransformer(512, 512)

	repository := NewUserRepository(database)
	service := NewUserService(repository, bucket, imageValidator, avatarResizer)
	handler := NewUserHandler(service)

	return UserModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module UserModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/users",
		func(router chi.Router) {
			router.Post("/", module.Handler.UserCreate)
			router.Get("/", module.Handler.GetUsers)
			router.Get("/{id}", module.Handler.GetUserById)
			router.Patch("/{id}", module.Handler.UpdateUserById)
			router.Patch("/{id}/avatar", module.Handler.UpdateUserAvatar)
			router.Patch("/{id}/role", module.Handler.AssignRole)
			router.Delete("/{id}", module.Handler.DeleteUser)
			router.Get("/me", module.Handler.GetSelf)
			router.Patch("/me/avatar", module.Handler.UpdateSelfAvatar)
		})
}
