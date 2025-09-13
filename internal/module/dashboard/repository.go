package dashboard

import (
	"bloggo/internal/module/dashboard/models"
	"database/sql"
	"os"
	"path/filepath"
	"syscall"
)

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) DashboardRepository {
	return DashboardRepository{
		db: db,
	}
}

func (repo *DashboardRepository) GetDashboardStats() (*models.ResponseDashboardStats, error) {
	stats := &models.ResponseDashboardStats{}

	// Get pending versions
	pendingVersions, err := repo.getPendingVersions()
	if err != nil {
		return nil, err
	}
	stats.PendingVersions = pendingVersions

	// Get recent activity
	recentActivity, err := repo.getRecentActivity()
	if err != nil {
		return nil, err
	}
	stats.RecentActivity = recentActivity

	// Get publishing rate
	publishingRate, err := repo.getPublishingRate()
	if err != nil {
		return nil, err
	}
	stats.PublishingRate = publishingRate

	// Get author performance
	authorPerformance, err := repo.getAuthorPerformance()
	if err != nil {
		return nil, err
	}
	stats.AuthorPerformance = authorPerformance

	// Get draft count
	draftCount, err := repo.getDraftCount()
	if err != nil {
		return nil, err
	}
	stats.DraftCount = draftCount

	// Get popular tags
	popularTags, err := repo.getPopularTags()
	if err != nil {
		return nil, err
	}
	stats.PopularTags = popularTags

	// Get storage usage
	storageUsage, err := repo.getStorageUsage()
	if err != nil {
		return nil, err
	}
	stats.StorageUsage = storageUsage

	return stats, nil
}

func (repo *DashboardRepository) getPendingVersions() ([]models.PendingVersion, error) {
	rows, err := repo.db.Query(QueryGetPendingVersions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.PendingVersion
	for rows.Next() {
		var version models.PendingVersion
		err := rows.Scan(&version.Id, &version.Title, &version.AuthorId, &version.AuthorName, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (repo *DashboardRepository) getRecentActivity() ([]models.RecentActivity, error) {
	rows, err := repo.db.Query(QueryGetRecentActivity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.RecentActivity
	for rows.Next() {
		var activity models.RecentActivity
		err := rows.Scan(&activity.Id, &activity.Title, &activity.PublishedAt)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (repo *DashboardRepository) getPublishingRate() (models.PublishingRate, error) {
	var rate models.PublishingRate

	err := repo.db.QueryRow(QueryGetPublishingRateWeek).Scan(&rate.ThisWeek)
	if err != nil {
		return rate, err
	}

	err = repo.db.QueryRow(QueryGetPublishingRateMonth).Scan(&rate.ThisMonth)
	if err != nil {
		return rate, err
	}

	return rate, nil
}

func (repo *DashboardRepository) getAuthorPerformance() ([]models.AuthorPerformance, error) {
	rows, err := repo.db.Query(QueryGetAuthorPerformance)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []models.AuthorPerformance
	for rows.Next() {
		var performance models.AuthorPerformance
		err := rows.Scan(&performance.AuthorId, &performance.AuthorName, &performance.PostCount)
		if err != nil {
			return nil, err
		}
		performances = append(performances, performance)
	}
	return performances, nil
}

func (repo *DashboardRepository) getDraftCount() (models.DraftCount, error) {
	var draftCount models.DraftCount

	err := repo.db.QueryRow(QueryGetTotalDraftCount).Scan(&draftCount.TotalDrafts)
	if err != nil {
		return draftCount, err
	}

	rows, err := repo.db.Query(QueryGetDraftsByAuthor)
	if err != nil {
		return draftCount, err
	}
	defer rows.Close()

	var draftsByAuthor []models.DraftsByAuthor
	for rows.Next() {
		var draft models.DraftsByAuthor
		err := rows.Scan(&draft.AuthorId, &draft.AuthorName, &draft.DraftCount)
		if err != nil {
			return draftCount, err
		}
		draftsByAuthor = append(draftsByAuthor, draft)
	}
	draftCount.DraftsByAuthor = draftsByAuthor

	return draftCount, nil
}

func (repo *DashboardRepository) getPopularTags() ([]models.PopularTag, error) {
	rows, err := repo.db.Query(QueryGetPopularTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.PopularTag
	for rows.Next() {
		var tag models.PopularTag
		err := rows.Scan(&tag.Id, &tag.Name, &tag.Usage)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (repo *DashboardRepository) getStorageUsage() (models.StorageUsage, error) {
	var storage models.StorageUsage

	// Get current working directory to calculate filesystem stats
	cwd, err := os.Getwd()
	if err != nil {
		return storage, err
	}

	// Get filesystem stats for the current directory
	var stat syscall.Statfs_t
	err = syscall.Statfs(cwd, &stat)
	if err != nil {
		return storage, err
	}

	// Calculate filesystem usage
	blockSize := uint64(stat.Bsize)
	totalBytes := stat.Blocks * blockSize
	freeBytes := stat.Bavail * blockSize
	usedBytes := totalBytes - freeBytes

	storage.FilesystemUsedBytes = int64(usedBytes)
	storage.FilesystemFreeBytes = int64(freeBytes)

	// Calculate bloggo storage usage from the uploads directory
	uploadDir := "uploads"
	bloggoUsed := int64(0)
	fileCount := 0

	err = filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If uploads directory doesn't exist, return zero values
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() {
			bloggoUsed += info.Size()
			fileCount++
		}
		return nil
	})

	if err != nil {
		// If there's an error (like directory not found), set zero values for bloggo usage
		bloggoUsed = 0
		fileCount = 0
	}

	storage.BloggoUsedBytes = bloggoUsed
	storage.FileCount = fileCount

	return storage, nil
}
