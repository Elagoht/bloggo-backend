package app

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"sync"

	"bloggo/internal/config"

	"github.com/go-chi/chi"
)

type Application struct {
	Database sql.DB
	Config   config.Config
	Router   chi.Mux
}

var (
	once     sync.Once
	instance Application
)

func GetInstance() *Application {
	once.Do(func() {
		instance = Application{
			Config: config.Load("bloggo-config.json"),
		}
	})
	return &instance
}

func (app *Application) Bootstrap() {
	portString := strconv.Itoa(app.Config.Port)
	// Start the server
	server := &http.Server{
		Addr:    ":" + portString,
		Handler: &app.Router,
	}
	log.Printf("Starting server on http://localhost:%s", portString)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
