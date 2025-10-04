package post

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"fmt"

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
		panic(fmt.Sprintf("Post module failed to create file storage: %v", err))
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

	router.Route("/posts", func(router chi.Router) {
		// Most routes require authentication
		router.With(middleware.AuthMiddleware(&config)).Group(func(router chi.Router) {
			router.Get("/", module.Handler.ListPosts)
			router.Get("/{id}", module.Handler.GetPostById)
			router.Post("/", module.Handler.CreatePostWithFirstVersion)
			router.Delete("/{id}", module.Handler.DeletePostById)
			router.Get("/{id}/versions", module.Handler.ListPostVersionsGetByPostId)
			router.Get("/{id}/versions/{versionId}", module.Handler.GetPostVersionById)
			router.Post("/{id}/versions", module.Handler.CreateVersionFromLatest)
			router.Post("/versions/{versionId}/duplicate", module.Handler.CreateVersionFromSpecificVersion)
			router.Patch("/{id}/versions/{versionId}", module.Handler.UpdateUnsubmittedOwnVersion)
			router.Post("/{id}/versions/{versionId}/submit", module.Handler.SubmitVersionForReview)
			router.Post("/{id}/versions/{versionId}/approve", module.Handler.ApproveVersion)
			router.Post("/{id}/versions/{versionId}/reject", module.Handler.RejectVersion)
			router.Post("/{id}/versions/{versionId}/publish", module.Handler.PublishVersion)
			router.Patch("/{id}/versions/{versionId}/category", module.Handler.UpdateVersionCategory)
			router.Delete("/{id}/versions/{versionId}", module.Handler.DeleteVersionById)
			router.Get("/{id}/versions/{versionId}/generative-fill", module.Handler.GenerativeFill)
			router.Post("/{id}/tags", module.Handler.AssignTagsToPost)
		})

		// Track-view endpoint only requires trusted frontend header
		router.With(
			middleware.TrustedFrontendMiddleware(&config),
		).Post("/track-view", module.Handler.TrackView)
	})
}
