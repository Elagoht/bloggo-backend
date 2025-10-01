package tags

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type TagsAPIModule struct {
	Handler    TagsAPIHandler
	Service    TagsAPIService
	Repository TagsAPIRepository
}

func NewModule() TagsAPIModule {
	database := db.Get()

	repository := NewTagsAPIRepository(database)
	service := NewTagsAPIService(repository)
	handler := NewTagsAPIHandler(service)

	return TagsAPIModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module TagsAPIModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.Route("/api/tags", func(r chi.Router) {
		// All endpoints require trusted frontend header
		r.Use(middleware.TrustedFrontendMiddleware(&config))

		r.Get("/", module.Handler.ListTags)
		r.Get("/{slug}", module.Handler.GetTagBySlug)
	})
}
