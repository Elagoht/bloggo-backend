package audit

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
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
	config := config.Get()

	// Create layers
	repository := NewAuditRepository(database)
	service := NewAuditService(repository)
	handler := NewAuditHandler(service)

	// Define routes
	router.With(middleware.AuthMiddleware(&config)).Route("/audit-logs", func(r chi.Router) {
		r.Get("/", handler.GetAuditLogs)
		r.Get("/entity/{type}/{id}", handler.GetAuditLogsByEntity)
		r.Get("/user/{id}", handler.GetAuditLogsByUser)
	})
}