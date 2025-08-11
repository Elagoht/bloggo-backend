package tag

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"bloggo/internal/utils/permission"

	"github.com/go-chi/chi"
)

type TagModule struct {
	Handler    TagHandler
	Service    TagService
	Repository TagRepository
}

func NewModule() TagModule {
	database := db.Get()
	repository := NewTagRepository(database)
	service := NewTagService(repository)
	handler := NewTagHandler(service)

	return TagModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module TagModule) RegisterModule(router *chi.Mux) {
	config := config.Get()
	permissionStore := permissions.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/tags",
		func(router chi.Router) {
			authority := permission.NewChecker(permissionStore)

			router.Get("/", module.Handler.GetCategories)
			router.Get("/{slug}", authority.Require("tag:manage", module.Handler.GetTagBySlug))
			router.Post("/", authority.Require("tag:manage", module.Handler.TagCreate))
			router.Patch("/{slug}", authority.Require("tag:manage", module.Handler.TagUpdate))
			router.Delete("/{slug}", authority.Require("tag:manage", module.Handler.TagDelete))
		},
	)
}
