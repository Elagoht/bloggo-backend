package dashboard

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/dashboard/models"
	"bloggo/internal/utils/apierrors"
)

type DashboardService struct {
	repository  DashboardRepository
	permissions permissions.Store
}

func NewDashboardService(repository DashboardRepository, permissions permissions.Store) DashboardService {
	return DashboardService{
		repository,
		permissions,
	}
}

func (service *DashboardService) GetDashboardStats(userRoleId int64) (*models.ResponseDashboardStats, error) {
	// Check permissions - user needs at least one of these permissions to see dashboard stats
	hasPostPublish := service.permissions.HasPermission(userRoleId, "post:publish")
	hasStatsViewOthers := service.permissions.HasPermission(userRoleId, "statistics:view-others")
	hasStatsViewTotal := service.permissions.HasPermission(userRoleId, "statistics:view-total")
	hasStatsViewSelf := service.permissions.HasPermission(userRoleId, "statistics:view-self")
	hasUserList := service.permissions.HasPermission(userRoleId, "user:list")
	hasTagList := service.permissions.HasPermission(userRoleId, "tag:list")
	hasTagView := service.permissions.HasPermission(userRoleId, "tag:view")
	hasPostCreate := service.permissions.HasPermission(userRoleId, "post:create")

	if !hasPostPublish && !hasStatsViewOthers && !hasStatsViewTotal && !hasStatsViewSelf && !hasUserList && !hasTagList && !hasTagView && !hasPostCreate {
		return nil, apierrors.ErrForbidden
	}

	// Get all stats first
	stats, err := service.repository.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	// Filter stats based on permissions
	filteredStats := &models.ResponseDashboardStats{}

	// Pending versions - only for users with post:publish
	if hasPostPublish {
		filteredStats.PendingVersions = stats.PendingVersions
	}

	// Recent activity - for users with statistics permissions
	if hasStatsViewOthers || hasStatsViewTotal || hasStatsViewSelf {
		filteredStats.RecentActivity = stats.RecentActivity
	}

	// Publishing rate - for users with statistics permissions
	if hasStatsViewOthers || hasStatsViewTotal || hasStatsViewSelf {
		filteredStats.PublishingRate = stats.PublishingRate
	}

	// Author performance - for users with statistics or user permissions
	if hasStatsViewOthers || hasStatsViewTotal || hasUserList {
		filteredStats.AuthorPerformance = stats.AuthorPerformance
	}

	// Draft count - for users with statistics:view-self or post:create
	if hasStatsViewSelf || hasPostCreate {
		filteredStats.DraftCount = stats.DraftCount
	}

	// Popular tags - for users with tag permissions
	if hasTagList || hasTagView {
		filteredStats.PopularTags = stats.PopularTags
	}

	// Storage usage - only for users with statistics:view-total or user:list (admin level)
	if hasStatsViewTotal || hasUserList {
		filteredStats.StorageUsage = stats.StorageUsage
	}

	return filteredStats, nil
}
