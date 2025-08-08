package tag

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"bloggo/internal/utils/permission"
	"database/sql"

	"github.com/go-chi/chi"
)

type TagModule struct {
	Handler     TagHandler
	Service     TagService
	Repository  TagRepository
	Config      *config.Config
	Permissions permissions.Store
}

func NewModule(
	database *sql.DB,
	config *config.Config,
	permissions permissions.Store,
) TagModule {
	repository := NewTagRepository(database)
	service := NewTagService(repository)
	handler := NewTagHandler(service)

	return TagModule{
		Handler:     handler,
		Service:     service,
		Repository:  repository,
		Config:      config,
		Permissions: permissions,
	}
}

func (module TagModule) RegisterModule(router *chi.Mux) {
	router.With(middleware.AuthMiddleware(module.Config)).Route(
		"/tags",
		func(router chi.Router) {
			authority := permission.NewChecker(module.Permissions)

			router.Get("/", module.Handler.GetCategories)
			router.Get("/{slug}", authority.Require("tag:manage", module.Handler.GetTagBySlug))
			router.Post("/", authority.Require("tag:manage", module.Handler.TagCreate))
			router.Patch("/{slug}", authority.Require("tag:manage", module.Handler.TagUpdate))
			router.Delete("/{slug}", authority.Require("tag:manage", module.Handler.TagDelete))
		},
	)
}
