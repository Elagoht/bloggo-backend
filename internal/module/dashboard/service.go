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
	hasUserList := service.permissions.HasPermission(userRoleId, "user:list")
	
	if !hasPostPublish && !hasStatsViewOthers && !hasStatsViewTotal && !hasUserList {
		return nil, apierrors.ErrForbidden
	}

	return service.repository.GetDashboardStats()
}