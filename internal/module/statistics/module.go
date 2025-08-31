package statistics

import (
	"bloggo/internal/config"
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/middleware"

	"github.com/go-chi/chi"
)

type StatisticsModule struct {
	Handler    StatisticsHandler
	Service    StatisticsService
	Repository StatisticsRepository
}

func NewModule() StatisticsModule {
	database := db.Get()
	permissionStore := permissions.Get()
	repository := NewStatisticsRepository(database)
	service := NewStatisticsService(repository, permissionStore)
	handler := NewStatisticsHandler(service)

	return StatisticsModule{
		Handler:    handler,
		Service:    service,
		Repository: repository,
	}
}

func (module StatisticsModule) RegisterModule(router *chi.Mux) {
	config := config.Get()

	router.With(middleware.AuthMiddleware(&config)).Route(
		"/statistics",
		func(router chi.Router) {
			// All statistics endpoint (for users with statistics:view-total permission)
			router.Get("/", module.Handler.GetAllStatistics)

			// User own statistics endpoint (for users with statistics:view-self permission)
			router.Get("/me", module.Handler.GetUserOwnStatistics)

			// Author statistics endpoint (for users with statistics:view-others permission)
			router.Get("/author/{authorId}", module.Handler.GetAuthorStatistics)
		},
	)
}
