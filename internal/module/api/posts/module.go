package posts

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type PostsAPIModule struct {
	Handler    PostsAPIHandler
	Service    PostsAPIService
	Repository PostsAPIRepository
}

func NewModule() PostsAPIModule {
	database := db.Get()

	repository := NewPostsAPIRepository(database)
	service := NewPostsAPIService(repository)
	handler := NewPostsAPIHandler(service)

	return PostsAPIModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module PostsAPIModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.Route("/api/posts", func(r chi.Router) {
		// All endpoints require trusted frontend header
		r.Use(middleware.TrustedFrontendMiddleware(&config))

		r.Get("/", module.Handler.ListPublishedPosts)
		r.Get("/{slug}", module.Handler.GetPublishedPostBySlug)
		r.Post("/{slug}/view", module.Handler.TrackPostView)
	})
}
