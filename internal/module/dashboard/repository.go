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

func (repo *DashboardRepository) GetPendingVersions() ([]models.PendingVersion, error) {
	rows, err := repo.db.Query(QueryGetPendingVersions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.PendingVersion
	for rows.Next() {
		var version models.PendingVersion
		err := rows.Scan(&version.Id, &version.PostId, &version.Title, &version.AuthorId, &version.AuthorName, &version.AuthorAvatar, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (repo *DashboardRepository) GetRecentActivity() ([]models.RecentActivity, error) {
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

func (repo *DashboardRepository) GetPublishingRate() (models.PublishingRate, error) {
	var rate models.PublishingRate

	err := repo.db.QueryRow(QueryGetPublishingRateDay).Scan(&rate.Today)
	if err != nil {
		return rate, err
	}

	err = repo.db.QueryRow(QueryGetPublishingRateWeek).Scan(&rate.ThisWeek)
	if err != nil {
		return rate, err
	}

	err = repo.db.QueryRow(QueryGetPublishingRateMonth).Scan(&rate.ThisMonth)
	if err != nil {
		return rate, err
	}

	err = repo.db.QueryRow(QueryGetPublishingRateYear).Scan(&rate.ThisYear)
	if err != nil {
		return rate, err
	}

	return rate, nil
}

func (repo *DashboardRepository) GetAuthorPerformance() ([]models.AuthorPerformance, error) {
	rows, err := repo.db.Query(QueryGetAuthorPerformance)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []models.AuthorPerformance
	for rows.Next() {
		var performance models.AuthorPerformance
		err := rows.Scan(&performance.Id, &performance.Name, &performance.Avatar, &performance.PostCount)
		if err != nil {
			return nil, err
		}
		performances = append(performances, performance)
	}
	return performances, nil
}

func (repo *DashboardRepository) GetDraftCount() (models.DraftCount, error) {
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
		err := rows.Scan(&draft.AuthorId, &draft.AuthorName, &draft.AuthorAvatar, &draft.DraftCount)
		if err != nil {
			return draftCount, err
		}
		draftsByAuthor = append(draftsByAuthor, draft)
	}
	draftCount.DraftsByAuthor = draftsByAuthor

	return draftCount, nil
}

func (repo *DashboardRepository) GetPopularTags() ([]models.PopularTag, error) {
	rows, err := repo.db.Query(QueryGetPopularTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.PopularTag
	for rows.Next() {
		var tag models.PopularTag
		err := rows.Scan(&tag.Id, &tag.Name, &tag.Slug, &tag.Usage)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (repo *DashboardRepository) GetStorageUsage() (models.StorageUsage, error) {
	var storage models.StorageUsage

	// Get filesystem stats for the root directory to get total system storage
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return storage, err
	}

	// Calculate filesystem usage
	blockSize := uint64(stat.Bsize)
	totalBytes := stat.Blocks * blockSize
	freeBytes := stat.Bavail * blockSize
	totalUsedBytes := totalBytes - freeBytes

	// Will be converted to decimal later

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

	// Convert binary bytes to decimal bytes (true GB/MB/KB)
	// Multiply by 1024^3 and divide by 1000^3 to convert GiB to GB
	// This makes filesystem stats match what users see in Finder/Explorer
	conversionFactor := 1.073741824 // 1024^3 / 1000^3

	storage.UsedByBloggoBytes = int64(float64(bloggoUsed) * conversionFactor)
	storage.UsedByOthersBytes = int64(float64(int64(totalUsedBytes)-bloggoUsed) * conversionFactor)
	storage.FreeBytes = int64(float64(freeBytes) * conversionFactor)
	storage.FileCount = fileCount

	return storage, nil
}
