package storage

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"database/sql"

	"github.com/go-chi/chi"
)

type StorageModule struct {
	Handler     StorageHandler
	Config      *config.Config
	Permissions permissions.Store
}

func NewModule(
	database *sql.DB,
	config *config.Config,
	permissions permissions.Store,
) StorageModule {
	handler := NewStorageHandler()

	return StorageModule{
		Handler:     handler,
		Config:      config,
		Permissions: permissions,
	}
}

func (module StorageModule) RegisterModule(router *chi.Mux) {
	router.Route("/uploads", func(router chi.Router) {
		router.With(middleware.AuthMiddleware(module.Config)).Route(
			"/uploads",
			func(router chi.Router) {
				// authority := permission.NewChecker(module.Permissions)
			},
		)
		// Public
		router.Get("/avatar/{imageId}", module.Handler.ServeUserAvatars)
	})
}
