package app

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module"
	"bloggo/internal/utils/audit"

	"github.com/go-chi/chi"
)

type Application struct {
	Router *chi.Mux
}

var (
	once     sync.Once
	instance Application
)

// Get singleton instance
func Get() *Application {
	once.Do(func() {
		// Initialize singletons
		databaseConnection := db.Get()
		permissionStore := permissions.Get()
		// Initial cache from database
		permissionStore.Load(databaseConnection)
		
		// Initialize audit logger
		audit.InitializeAuditLogger(databaseConnection)

		instance = Application{
			Router: chi.NewRouter(),
		}
	})
	return &instance
}

func (app *Application) RegisterModules(modules []module.Module) {
	for _, module := range modules {
		module.RegisterModule(app.Router)
	}
}

func (app *Application) RegisterGlobalMiddlewares(
	middlewares []func(http.Handler) http.Handler,
) {
	for _, middleware := range middlewares {
		app.Router.Use(middleware)
	}
}

func (app *Application) Bootstrap() {
	config := config.Get()
	portString := strconv.Itoa(config.Port)
	// Start the server
	server := &http.Server{
		Addr:    ":" + portString,
		Handler: app.Router,
	}
	log.Printf("Starting server on http://localhost:%s", portString)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
