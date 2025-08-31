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
	hasViewAllPermission := service.permissions.HasPermission(userRoleId, "statistics:view-total")

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

	// Get category length distribution
	categoryLength, err := service.repository.GetCategoryLengthDistribution()
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
		ViewStats:                     viewStats,
		Last24Hours:                   last24Hours,
		CategoryViewsDistribution:     categoryViews,
		MostViewedBlogs:               mostViewed,
		BlogStats:                     blogStats,
		LongestBlogs:                  longestBlogs,
		CategoryBlogsDistribution:     categoryBlogs,
		CategoryReadTimeDistribution:  categoryReadTime,
		CategoryLengthDistribution:    categoryLength,
		TopUserAgents:                 topUserAgents,
		DeviceTypeDistribution:        deviceDistribution,
		OSDistribution:                osDistribution,
		BrowserDistribution:           browserDistribution,
	}, nil
}

func (service *StatisticsService) GetUserOwnStatistics(userRoleId int64, userId int64) (*models.ResponseAuthorStatistics, error) {
	hasViewSelfPermission := service.permissions.HasPermission(userRoleId, "statistics:view-self")
	if !hasViewSelfPermission {
		return nil, apierrors.ErrForbidden
	}

	return service.GetAuthorStatistics(userId, userRoleId, userId)
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
		// If viewing another user's stats, need statistics:view-others permission
		hasViewOthersPermission := service.permissions.HasPermission(userRoleId, "statistics:view-others")
		if !hasViewOthersPermission {
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

	// Get category length distribution (author-specific)
	categoryLength, err := service.repository.GetCategoryLengthDistribution()
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
		AuthorStats:                   authorStats,
		ViewStats:                     viewStats,
		Last24Hours:                   last24Hours,
		CategoryViewsDistribution:     categoryViews,
		MostViewedBlogs:               mostViewed,
		BlogStats:                     blogStats,
		LongestBlogs:                  longestBlogs,
		CategoryBlogsDistribution:     categoryBlogs,
		CategoryReadTimeDistribution:  categoryReadTime,
		CategoryLengthDistribution:    categoryLength,
		TopUserAgents:                 topUserAgents,
		DeviceTypeDistribution:        deviceDistribution,
		OSDistribution:                osDistribution,
		BrowserDistribution:           browserDistribution,
	}, nil
}

