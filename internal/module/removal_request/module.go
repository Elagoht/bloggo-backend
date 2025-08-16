package removal_request

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"log"

	"github.com/go-chi/chi"
)

type RemovalRequestModule struct {
	Handler    RemovalRequestHandler
	Service    RemovalRequestService
	Repository RemovalRequestRepository
}

func NewModule() RemovalRequestModule {
	database := db.Get()
	permissionStore := permissions.Get()
	
	// Use the same bucket as posts for cover image cleanup
	bucketInstance, err := bucket.NewFileSystemBucket("posts/versions/covers")
	if err != nil {
		log.Fatalln("Removal request module cannot create file storage")
	}

	repository := NewRemovalRequestRepository(database)
	service := NewRemovalRequestService(repository, permissionStore, bucketInstance)
	handler := NewRemovalRequestHandler(service)

	return RemovalRequestModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module RemovalRequestModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/removal-requests",
		func(router chi.Router) {
			// Create removal request
			router.Post("/", module.Handler.CreateRemovalRequest)
			
			// Get all removal requests (admin/editor only)
			router.Get("/", module.Handler.GetRemovalRequestList)
			
			// Get user's own removal requests
			router.Get("/mine", module.Handler.GetUserRemovalRequests)
			
			// Get specific removal request
			router.Get("/{id}", module.Handler.GetRemovalRequestById)
			
			// Approve removal request (admin/editor only)
			router.Post("/{id}/approve", module.Handler.ApproveRemovalRequest)
			
			// Reject removal request (admin/editor only)
			router.Post("/{id}/reject", module.Handler.RejectRemovalRequest)
		},
	)
}