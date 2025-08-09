package post

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"database/sql"

	"github.com/go-chi/chi"
)

type PostModule struct {
	Handler     PostHandler
	Service     PostService
	Repository  PostRepository
	Config      *config.Config
	Permissions permissions.Store
}

func NewModule(
	database *sql.DB,
	config *config.Config,
	permissions permissions.Store,
) PostModule {
	repository := NewPostRepository(database)
	service := NewPostService(repository)
	handler := NewPostHandler(service)

	return PostModule{
		Handler:     handler,
		Service:     service,
		Repository:  repository,
		Config:      config,
		Permissions: permissions,
	}
}

func (module PostModule) RegisterModule(router *chi.Mux) {
	router.With(middleware.AuthMiddleware(module.Config)).Route(
		"/posts",
		func(router chi.Router) {
			// authority := permission.NewChecker(module.Permissions)

			router.Get("/", module.Handler.ListPosts)
			router.Get("/{slug}", module.Handler.GetPostBySlug)
			router.Post("/", module.Handler.CreatePostWithFirstVersion)
		},
	)
}
