package main

import (
	"bloggo/internal/app"
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
	"bloggo/internal/module/statistics"
	"bloggo/internal/module/storage"
	"bloggo/internal/module/tag"
	"bloggo/internal/module/user"
	"net/http"
)

func main() {
	// Get singleton application
	application := app.Get()

	// Register global middlewares
	middlewares := []func(http.Handler) http.Handler{
		middleware.ResponseJSON,
		middleware.AllowSpecificOrigin,
		middleware.GlobalRateLimiter(),
	}
	application.RegisterGlobalMiddlewares(middlewares)

	// Register modules
	modules := []module.Module{
		api.NewModule(), // Public API module (must be first)
		category.NewModule(),
		tag.NewModule(),
		post.NewModule(),
		user.NewModule(),
		session.NewModule(),
		storage.NewModule(),
		removal_request.NewModule(),
		statistics.NewModule(),
		audit.NewModule(),
		dashboard.NewModule(),
		search.NewModule(),
		health.NewModule(),
		keyvalue.NewModule(),
	}
	application.RegisterModules(modules)

	// Start app
	application.Bootstrap()
}
