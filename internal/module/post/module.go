package post

import (
	"bloggo/internal/config"
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"database/sql"
	"log"

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
	bucket, err := bucket.NewFileSystemBucket("posts/versions/covers")
	if err != nil {
		log.Fatalln("Post module cannot created file storage")
	}
	imageValidator := validatefile.NewImageValidator(10 << 20) // 5MB
	coverResizer := transformfile.NewImageTransformer(1280, 720)

	repository := NewPostRepository(database)
	service := NewPostService(repository, bucket, imageValidator, coverResizer)
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
			router.Get("/{id}", module.Handler.GetPostById)
			router.Post("/", module.Handler.CreatePostWithFirstVersion)
		},
	)
}
