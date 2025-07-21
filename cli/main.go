package main

import (
	"bloggo/internal/app"
	"bloggo/internal/middleware"
	"bloggo/internal/module"
	"bloggo/internal/module/auth"
	"bloggo/internal/module/category"
	"bloggo/internal/module/user"
	"net/http"
)

func main() {
	// Create singleton application
	application := app.GetInstance()

	// Register global middlewares
	middlewares := []func(http.Handler) http.Handler{
		middleware.ResponseJSON,
	}
	application.RegisterGlobalMiddlewares(middlewares)

	// Register modules
	modules := []module.Module{
		category.NewModule(application.Database, &application.Config, application.Permissions),
		user.NewModule(application.Database),
		auth.NewModule(application.Database, &application.Config),
	}
	application.RegisterModules(modules)

	// Start app
	application.Bootstrap()
}
