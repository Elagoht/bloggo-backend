package statistics

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/statistics/models"
	"bloggo/internal/utils/apierrors"
)

type StatisticsService struct {
	repository  StatisticsRepository
	permissions permissions.Store
}

func NewStatisticsService(repository StatisticsRepository, permissions permissions.Store) StatisticsService {
	return StatisticsService{
		repository,
		permissions,
	}
}

func (service *StatisticsService) GetAllStatistics(userRoleId int64) (*models.ResponseAllStatistics, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	// Get view statistics
	viewStats, err := service.repository.GetViewStatistics()
	if err != nil {
		return nil, err
	}

	// Get last 24 hours views
	last24Hours, err := service.repository.GetLast24HoursViews()
	if err != nil {
		return nil, err
	}

	// Get category views distribution
	categoryViews, err := service.repository.GetCategoryViewsDistribution()
	if err != nil {
		return nil, err
	}

	// Get most viewed blogs (top 10)
	mostViewed, err := service.repository.GetMostViewedBlogs(10)
	if err != nil {
		return nil, err
	}

	// Get blog statistics
	blogStats, err := service.repository.GetBlogStatistics()
	if err != nil {
		return nil, err
	}

	// Get longest blogs (top 10)
	longestBlogs, err := service.repository.GetLongestBlogs(10)
	if err != nil {
		return nil, err
	}

	// Get category blog distribution
	categoryBlogs, err := service.repository.GetCategoryBlogDistribution()
	if err != nil {
		return nil, err
	}

	// Get category read time distribution
	categoryReadTime, err := service.repository.GetCategoryReadTimeDistribution()
	if err != nil {
		return nil, err
	}

	// Get top user agents (top 10)
	topUserAgents, err := service.repository.GetTopUserAgents(10)
	if err != nil {
		return nil, err
	}

	// Get device type distribution
	deviceDistribution, err := service.repository.GetDeviceTypeDistribution()
	if err != nil {
		return nil, err
	}

	// Get OS distribution
	osDistribution, err := service.repository.GetOSDistribution()
	if err != nil {
		return nil, err
	}

	// Get browser distribution
	browserDistribution, err := service.repository.GetBrowserDistribution()
	if err != nil {
		return nil, err
	}

	return &models.ResponseAllStatistics{
		ViewStats:                    viewStats,
		Last24Hours:                  last24Hours,
		CategoryViewsDistribution:    categoryViews,
		MostViewedBlogs:              mostViewed,
		BlogStats:                    blogStats,
		LongestBlogs:                 longestBlogs,
		CategoryBlogsDistribution:    categoryBlogs,
		CategoryReadTimeDistribution: categoryReadTime,
		TopUserAgents:                topUserAgents,
		DeviceTypeDistribution:       deviceDistribution,
		OSDistribution:               osDistribution,
		BrowserDistribution:          browserDistribution,
	}, nil
}

func (service *StatisticsService) GetAuthorStatistics(authorId int64, userRoleId int64, requestingUserId int64) (*models.ResponseAuthorStatistics, error) {
	// Check if requesting own statistics
	isOwnStats := requestingUserId == authorId

	// If viewing own stats, just need statistics:view-self permission
	if isOwnStats {
		hasViewSelfPermission := service.permissions.HasPermission(userRoleId, "statistics:view-self")
		if !hasViewSelfPermission {
			return nil, apierrors.ErrForbidden
		}
	} else {
		// If viewing another user's stats, need statistics:view-all permission
		hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")
		if !hasViewAllPermission {
			return nil, apierrors.ErrForbidden
		}
	}

	// Get author statistics
	authorStats, err := service.repository.GetAuthorStatistics(authorId)
	if err != nil {
		return nil, err
	}

	// Get author view statistics
	viewStats, err := service.repository.GetAuthorViewStatistics(authorId)
	if err != nil {
		return nil, err
	}

	// Get last 24 hours views (this will be global for now, could be author-specific later)
	last24Hours, err := service.repository.GetLast24HoursViews()
	if err != nil {
		return nil, err
	}

	// Get author category views distribution
	categoryViews, err := service.repository.GetAuthorCategoryViewsDistribution(authorId)
	if err != nil {
		return nil, err
	}

	// Get author's most viewed blogs (top 10)
	mostViewed, err := service.repository.GetAuthorMostViewedBlogs(authorId, 10)
	if err != nil {
		return nil, err
	}

	// Get author blog statistics
	blogStats, err := service.repository.GetAuthorBlogStatistics(authorId)
	if err != nil {
		return nil, err
	}

	// Get author's longest blogs (top 10)
	longestBlogs, err := service.repository.GetAuthorLongestBlogs(authorId, 10)
	if err != nil {
		return nil, err
	}

	// Get category blog distribution (author-specific)
	categoryBlogs, err := service.repository.GetCategoryBlogDistribution()
	if err != nil {
		return nil, err
	}

	// Get category read time distribution (author-specific)
	categoryReadTime, err := service.repository.GetCategoryReadTimeDistribution()
	if err != nil {
		return nil, err
	}

	// Get top user agents (global for now)
	topUserAgents, err := service.repository.GetTopUserAgents(10)
	if err != nil {
		return nil, err
	}

	// Get device type distribution (global for now)
	deviceDistribution, err := service.repository.GetDeviceTypeDistribution()
	if err != nil {
		return nil, err
	}

	// Get OS distribution (global for now)
	osDistribution, err := service.repository.GetOSDistribution()
	if err != nil {
		return nil, err
	}

	// Get browser distribution (global for now)
	browserDistribution, err := service.repository.GetBrowserDistribution()
	if err != nil {
		return nil, err
	}

	return &models.ResponseAuthorStatistics{
		AuthorStats:                  authorStats,
		ViewStats:                    viewStats,
		Last24Hours:                  last24Hours,
		CategoryViewsDistribution:    categoryViews,
		MostViewedBlogs:              mostViewed,
		BlogStats:                    blogStats,
		LongestBlogs:                 longestBlogs,
		CategoryBlogsDistribution:    categoryBlogs,
		CategoryReadTimeDistribution: categoryReadTime,
		TopUserAgents:                topUserAgents,
		DeviceTypeDistribution:       deviceDistribution,
		OSDistribution:               osDistribution,
		BrowserDistribution:          browserDistribution,
	}, nil
}

func (service *StatisticsService) GetViewStatistics(userRoleId int64) (*models.ViewStatistics, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetViewStatistics()
}

func (service *StatisticsService) GetLast24HoursViews(userRoleId int64) (*models.Last24HoursViews, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetLast24HoursViews()
}

func (service *StatisticsService) GetCategoryViewsDistribution(userRoleId int64) ([]models.CategoryViewDistribution, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetCategoryViewsDistribution()
}

func (service *StatisticsService) GetMostViewedBlogs(limit int, userRoleId int64) ([]models.MostViewedBlog, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}

	return service.repository.GetMostViewedBlogs(limit)
}

func (service *StatisticsService) GetBlogStatistics(userRoleId int64) (*models.BlogStatistics, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetBlogStatistics()
}

func (service *StatisticsService) GetLongestBlogs(limit int, userRoleId int64) ([]models.LongestBlog, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}

	return service.repository.GetLongestBlogs(limit)
}

func (service *StatisticsService) GetCategoryBlogDistribution(userRoleId int64) ([]models.CategoryBlogDistribution, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetCategoryBlogDistribution()
}

func (service *StatisticsService) GetCategoryReadTimeDistribution(userRoleId int64) ([]models.CategoryReadTimeDistribution, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetCategoryReadTimeDistribution()
}

func (service *StatisticsService) GetTopUserAgents(limit int, userRoleId int64) ([]models.UserAgentStat, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	if limit <= 0 || limit > 100 {
		limit = 10 // Default limit
	}

	return service.repository.GetTopUserAgents(limit)
}

func (service *StatisticsService) GetDeviceTypeDistribution(userRoleId int64) ([]models.DeviceTypeStat, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetDeviceTypeDistribution()
}

func (service *StatisticsService) GetOSDistribution(userRoleId int64) ([]models.OSStatistic, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetOSDistribution()
}

func (service *StatisticsService) GetBrowserDistribution(userRoleId int64) ([]models.BrowserStat, error) {
	// Everyone can view all statistics now
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-all")

	if !hasViewAllPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetBrowserDistribution()
}
