package category

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"bloggo/internal/module/ai"

	"github.com/go-chi/chi"
)

type CategoryModule struct {
	Handler    CategoryHandler
	Service    CategoryService
	Repository CategoryRepository
}

func NewModule() CategoryModule {
	database := db.Get()
	permissionStore := permissions.Get()
	aiService := ai.NewAIService()
	repository := NewCategoryRepository(database)
	service := NewCategoryService(repository, permissionStore, aiService)
	handler := NewCategoryHandler(service)

	return CategoryModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module CategoryModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/categories",
		func(router chi.Router) {
			router.Get("/", module.Handler.GetCategories)
			router.Get("/list", module.Handler.GetCategoryList)
			router.Get("/{slug}", module.Handler.GetCategoryBySlug)
			router.Post("/", module.Handler.CategoryCreate)
			router.Patch("/{slug}", module.Handler.CategoryUpdate)
			router.Delete("/{slug}", module.Handler.CategoryDelete)
			router.Get("/generative-fill", module.Handler.GenerativeFill)
		},
	)
}
