package dashboard

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/dashboard/models"
	"fmt"
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

func (service *DashboardService) GetDashboardStats(
	userRoleId int64,
) (*models.ResponseDashboardStats, error) {
	result := &models.ResponseDashboardStats{}

	// Set public stats
	publishingRate, err := service.repository.GetPublishingRate()
	if err != nil {
		return nil, err
	}
	result.PublishingRate = publishingRate

	draftCount, err := service.repository.GetDraftCount()
	if err != nil {
		return nil, err
	}

	for index := range draftCount.DraftsByAuthor {
		if draftCount.DraftsByAuthor[index].AuthorAvatar != nil && *draftCount.DraftsByAuthor[index].AuthorAvatar != "" {
			avatarPath := fmt.Sprintf("/uploads/avatar/%s", *draftCount.DraftsByAuthor[index].AuthorAvatar)
			draftCount.DraftsByAuthor[index].AuthorAvatar = &avatarPath
		}
	}

	result.DraftCount = draftCount

	popularTags, err := service.repository.GetPopularTags()
	if err != nil {
		return nil, err
	}
	result.PopularTags = popularTags

	// Pending versions
	if service.permissions.HasPermission(userRoleId, "post:publish") {
		pendingVersions, err := service.repository.GetPendingVersions()
		if err != nil {
			return nil, err
		}

		for index := range pendingVersions {
			if pendingVersions[index].AuthorAvatar != nil && *pendingVersions[index].AuthorAvatar != "" {
				avatarPath := fmt.Sprintf("/uploads/avatar/%s", *pendingVersions[index].AuthorAvatar)
				pendingVersions[index].AuthorAvatar = &avatarPath
			}
		}

		result.PendingVersions = pendingVersions
	}

	// Audit Logs
	if service.permissions.HasPermission(userRoleId, "auditlog:view") {
		recentActivity, err := service.repository.GetRecentActivity()
		if err != nil {
			return nil, err
		}
		result.RecentActivity = recentActivity
	}

	// Performance and system details
	if service.permissions.HasPermission(userRoleId, "statistics:view-others") {
		authorPerformance, err := service.repository.GetAuthorPerformance()
		if err != nil {
			return nil, err
		}
		result.AuthorPerformance = authorPerformance

		for index := range authorPerformance {
			if authorPerformance[index].Avatar != nil && *authorPerformance[index].Avatar != "" {
				avatarPath := fmt.Sprintf("/uploads/avatar/%s", *authorPerformance[index].Avatar)
				authorPerformance[index].Avatar = &avatarPath
			}
		}

		storageUsage, err := service.repository.GetStorageUsage()
		if err != nil {
			return nil, err
		}
		result.StorageUsage = storageUsage
	}

	return result, nil
}
