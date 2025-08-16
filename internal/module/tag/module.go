package tag

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type TagModule struct {
	Handler    TagHandler
	Service    TagService
	Repository TagRepository
}

func NewModule() TagModule {
	database := db.Get()
	permissionStore := permissions.Get()

	repository := NewTagRepository(database)
	service := NewTagService(repository, permissionStore)
	handler := NewTagHandler(service)

	return TagModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module TagModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/tags",
		func(router chi.Router) {
			// Public routes
			router.Get("/", module.Handler.GetCategories)
			router.Get("/{slug}", module.Handler.GetTagBySlug)

			// Editor-only routes
			router.Post("/", module.Handler.TagCreate)
			router.Patch("/{slug}", module.Handler.TagUpdate)
			router.Delete("/{slug}", module.Handler.TagDelete)
		},
	)

	// Post-tag relationship routes
	router.With(middleware.AuthMiddleware(&config)).Route(
		"/posts/{postId}/tags",
		func(router chi.Router) {
			// Public route to view post tags
			router.Get("/", module.Handler.GetPostTags)

			// Editor-only routes for managing post tags
			router.Post("/", module.Handler.AssignTagsToPost)
			router.Delete("/{tagId}", module.Handler.RemoveTagFromPost)
		},
	)
}
