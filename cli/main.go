package main

import (
	"bloggo/internal/app"
	"bloggo/internal/module"
	"bloggo/internal/module/category"
)

func main() {
	// Create singleton application
	application := app.GetInstance()

	// Register odules
	modules := []module.Module{
		category.NewModule(application.Database),
	}
	application.RegisterModules(modules)

	// Start app
	application.Bootstrap()
}
