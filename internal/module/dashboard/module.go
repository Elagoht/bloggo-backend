package dashboard

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"
	"github.com/go-chi/chi"
)

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) RegisterModule(router *chi.Mux) {
	// Initialize dependencies
	database := db.Get()
	permissionStore := permissions.Get()
	config := config.Get()

	// Create layers
	repository := NewDashboardRepository(database)
	service := NewDashboardService(repository, permissionStore)
	handler := NewDashboardHandler(service)

	// Define routes
	router.With(middleware.AuthMiddleware(&config)).Route("/dashboard", func(r chi.Router) {
		r.Get("/stats", handler.GetDashboardStats)
	})
}