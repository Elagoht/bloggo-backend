package main

import (
	"bloggo/internal/app"
	"bloggo/internal/middleware"
	"bloggo/internal/module"
	"bloggo/internal/module/category"
	"bloggo/internal/module/health"
	"bloggo/internal/module/post"
	"bloggo/internal/module/session"
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
		category.NewModule(),
		tag.NewModule(),
		post.NewModule(),
		user.NewModule(),
		session.NewModule(),
		storage.NewModule(),
		health.NewModule(),
	}
	application.RegisterModules(modules)

	// Start app
	application.Bootstrap()
}
