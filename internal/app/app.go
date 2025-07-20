package app

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"sync"

	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/module"

	"github.com/go-chi/chi"
)

type Application struct {
	Database *sql.DB
	Config   config.Config
	Router   *chi.Mux
}

var (
	once     sync.Once
	instance Application
)

// Get singleton instance
func GetInstance() *Application {
	once.Do(func() {
		instance = Application{
			Database: db.GetInstance(),
			Config:   config.Load("bloggo-config.json"),
			Router:   chi.NewRouter(),
		}
	})
	return &instance
}

func (app *Application) RegisterModules(modules []module.Module) {
	for _, module := range modules {
		module.RegisterModule(app.Database, app.Router)
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
	portString := strconv.Itoa(app.Config.Port)
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
