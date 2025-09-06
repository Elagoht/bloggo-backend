package audit

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
	repository := NewAuditRepository(database)
	service := NewAuditService(repository, permissionStore)
	handler := NewAuditHandler(service)

	// Define routes
	router.With(middleware.AuthMiddleware(&config)).Route("/audit", func(r chi.Router) {
		r.Get("/logs", handler.GetAuditLogs)
	})
}