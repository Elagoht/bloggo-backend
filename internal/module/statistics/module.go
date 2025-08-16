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
			// All statistics endpoint (supports optional ?userId=X query parameter)
			router.Get("/", module.Handler.GetAllStatistics)

			// Individual statistics endpoints
			router.Get("/views", module.Handler.GetViewStatistics)
			router.Get("/views/last-24-hours", module.Handler.GetLast24HoursViews)
			router.Get("/blogs", module.Handler.GetBlogStatistics)
			router.Get("/blogs/most-viewed", module.Handler.GetMostViewedBlogs)
			router.Get("/blogs/longest", module.Handler.GetLongestBlogs)
			router.Get("/categories/views", module.Handler.GetCategoryViewsDistribution)
			router.Get("/categories/blogs", module.Handler.GetCategoryBlogDistribution)
			router.Get("/categories/read-time", module.Handler.GetCategoryReadTimeDistribution)
			router.Get("/user-agents", module.Handler.GetTopUserAgents)
			router.Get("/device-types", module.Handler.GetDeviceTypeDistribution)
			router.Get("/operating-systems", module.Handler.GetOSDistribution)
			router.Get("/browsers", module.Handler.GetBrowserDistribution)
		},
	)
}
