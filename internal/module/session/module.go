package session

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/tokens"
	"bloggo/internal/module/user"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"log"

	"github.com/go-chi/chi"
)

type SessionModule struct {
	Handler    SessionHandler
	Service    SessionService
	Repository SessionRepository
}

func NewModule() SessionModule {
	database := db.Get()
	config := config.Get()
	refreshStore := tokens.GetStore()

	// Create user service for updating last login
	bucket, err := bucket.NewFileSystemBucket("users/avatars")
	if err != nil {
		log.Fatalln("Session module cannot create file storage for user service")
	}
	imageValidator := validatefile.NewImageValidator(5 << 20) // 5MB
	avatarResizer := transformfile.NewImageTransformer(512, 512)
	userRepository := user.NewUserRepository(database)
	userService := user.NewUserService(userRepository, bucket, imageValidator, avatarResizer)

	repository := NewSessionRepository(database)
	service := NewSessionService(repository, &config, refreshStore, &userService)
	handler := NewSessionHandler(service, &config)

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
