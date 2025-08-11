package storage

import (
	"bloggo/internal/config"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type StorageModule struct {
	Handler StorageHandler
}

func NewModule() StorageModule {
	handler := NewStorageHandler()

	return StorageModule{
		Handler: handler,
	}
}

func (module StorageModule) RegisterModule(router *chi.Mux) {
	config := config.Get()
	
	router.Route("/uploads", func(router chi.Router) {
		router.With(middleware.AuthMiddleware(&config)).Route(
			"/uploads",
			func(router chi.Router) {
				// Add authenticated storage routes here if needed
			},
		)
		// Public
		router.Get("/avatar/{imageId}", module.Handler.ServeUserAvatars)
	})
}
