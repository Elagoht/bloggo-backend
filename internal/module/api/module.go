package api

import (
	"bloggo/internal/module"
	"bloggo/internal/module/api/authors"
	"bloggo/internal/module/api/categories"
	"bloggo/internal/module/api/posts"
	"bloggo/internal/module/api/tags"

	"github.com/go-chi/chi"
)

type APIModule struct {
	PostsModule      posts.PostsAPIModule
	CategoriesModule categories.CategoriesAPIModule
	TagsModule       tags.TagsAPIModule
	AuthorsModule    authors.AuthorsAPIModule
}

func NewModule() APIModule {
	postsModule := posts.NewModule()
	categoriesModule := categories.NewModule()
	tagsModule := tags.NewModule()
	authorsModule := authors.NewModule()

	return APIModule{
		PostsModule:      postsModule,
		CategoriesModule: categoriesModule,
		TagsModule:       tagsModule,
		AuthorsModule:    authorsModule,
	}
}

func (apiModule APIModule) RegisterModule(router *chi.Mux) {
	// Register all API sub-modules
	subModules := []module.Module{
		apiModule.PostsModule,
		apiModule.CategoriesModule,
		apiModule.TagsModule,
		apiModule.AuthorsModule,
	}

	for _, subModule := range subModules {
		subModule.RegisterModule(router)
	}
}
