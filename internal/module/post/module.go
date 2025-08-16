package post

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"log"

	"github.com/go-chi/chi"
)

type PostModule struct {
	Handler    PostHandler
	Service    PostService
	Repository PostRepository
}

func NewModule() PostModule {
	database := db.Get()
	bucket, err := bucket.NewFileSystemBucket("posts/versions/covers")
	if err != nil {
		log.Fatalln("Post module cannot created file storage")
	}
	imageValidator := validatefile.NewImageValidator(10 << 20) // 5MB
	coverResizer := transformfile.NewImageTransformer(1280, 720)
	permissionStore := permissions.Get()

	repository := NewPostRepository(database)
	service := NewPostService(repository, bucket, imageValidator, coverResizer, permissionStore)
	handler := NewPostHandler(service)

	return PostModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module PostModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/posts",
		func(router chi.Router) {
			router.Get("/", module.Handler.ListPosts)
			router.Get("/{id}", module.Handler.GetPostById)
			router.Post("/", module.Handler.CreatePostWithFirstVersion)
			router.Delete("/{id}", module.Handler.DeletePostById)
			router.Get("/{id}/versions", module.Handler.ListPostVersionsGetByPostId)
			router.Get("/{id}/versions/{versionId}", module.Handler.GetPostVersionById)
			router.Post("/{id}/versions", module.Handler.CreateVersionFromLatest)
			router.Patch("/{id}/versions/{versionId}", module.Handler.UpdateUnsubmittedOwnVersion)
			router.Post("/{id}/versions/{versionId}/submit", module.Handler.SubmitVersionForReview)
			router.Post("/{id}/versions/{versionId}/approve", module.Handler.ApproveVersion)
			router.Post("/{id}/versions/{versionId}/reject", module.Handler.RejectVersion)
			router.Delete("/{id}/versions/{versionId}", module.Handler.DeleteVersionById)
		},
	)
}
