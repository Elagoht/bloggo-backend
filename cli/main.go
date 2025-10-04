package main

import (
	"bloggo/internal/app"
	"bloggo/internal/config"
	"bloggo/internal/db"
	embedpkg "bloggo/internal/embed"
	"bloggo/internal/middleware"
	"bloggo/internal/module"
	"bloggo/internal/module/api"
	"bloggo/internal/module/audit"
	"bloggo/internal/module/category"
	"bloggo/internal/module/dashboard"
	"bloggo/internal/module/health"
	"bloggo/internal/module/keyvalue"
	"bloggo/internal/module/post"
	"bloggo/internal/module/removal_request"
	"bloggo/internal/module/search"
	"bloggo/internal/module/session"
	"bloggo/internal/module/static"
	"bloggo/internal/module/statistics"
	"bloggo/internal/module/storage"
	"bloggo/internal/module/tag"
	"bloggo/internal/module/user"
	"bloggo/internal/module/webhook"
	"bloggo/internal/utils/validate"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	// Initialize validator
	if err := validate.MustInitialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize validator: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	if err := config.MustLoad(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	if err := db.MustConnect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Get singleton application
	application := app.Get()

	// Load embedded frontend
	distFS, err := embedpkg.GetDistFS()
	if err != nil {
		log.Printf("Warning: Failed to load embedded frontend: %v", err)
		log.Println("Starting without frontend - only API routes will be available")
	} else {
		// Register static assets WITHOUT any middlewares
		staticModule := static.NewModule(distFS)
		application.Router.Get("/assets/*", staticModule.ServeAssets)
		application.Router.Get("/favicon.ico", staticModule.ServeStaticFile)
		application.Router.Get("/manifest.json", staticModule.ServeStaticFile)
		application.Router.Get("/robots.txt", staticModule.ServeStaticFile)
		application.Router.Get("/api-docs.json", staticModule.ServeStaticFile)
		log.Println("Static assets registered")
	}

	// Create middlewares to apply to API/internal routes only
	middlewares := []func(http.Handler) http.Handler{
		middleware.ResponseJSON,
		middleware.AllowSpecificOrigin,
		middleware.GlobalRateLimiter(),
	}

	// Register public API module with middlewares
	application.Router.Group(func(r chi.Router) {
		// Apply middlewares to this group
		for _, mw := range middlewares {
			r.Use(mw)
		}

		// Register API modules
		apiModules := []module.Module{
			api.NewModule(), // Public API routes at /api/*
		}

		for _, mod := range apiModules {
			mod.RegisterModule(r.(*chi.Mux))
		}
	})

	// Register public storage module (uploads must be publicly accessible)
	storage.NewModule().RegisterModule(application.Router)

	// Register internal panel modules under /internal prefix with middlewares
	application.Router.Group(func(r chi.Router) {
		// Apply middlewares to this group
		for _, mw := range middlewares {
			r.Use(mw)
		}

		internalRouter := chi.NewRouter()
		internalModules := []module.Module{
			category.NewModule(),
			tag.NewModule(),
			post.NewModule(),
			user.NewModule(),
			session.NewModule(),
			removal_request.NewModule(),
			statistics.NewModule(),
			audit.NewModule(),
			dashboard.NewModule(),
			search.NewModule(),
			health.NewModule(),
			keyvalue.NewModule(),
			webhook.NewModule(),
		}

		for _, mod := range internalModules {
			mod.RegisterModule(internalRouter)
		}

		// Mount the internal router at /internal
		r.Mount("/internal", internalRouter)
	})

	// Register SPA catch-all (must be last)
	if distFS != nil {
		staticModule := static.NewModule(distFS)
		application.Router.NotFound(staticModule.ServeSPA)
		log.Println("SPA routing configured")
	}

	// Start app
	if err := application.Bootstrap(); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed to start: %v\n", err)
		os.Exit(1)
	}
}
