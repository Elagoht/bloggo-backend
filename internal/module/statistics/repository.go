package statistics

import (
	"bloggo/internal/module/statistics/models"
	"bloggo/internal/utils/useragent"
	"database/sql"
)

type StatisticsRepository struct {
	database *sql.DB
}

func NewStatisticsRepository(database *sql.DB) StatisticsRepository {
	return StatisticsRepository{
		database,
	}
}

func (repository *StatisticsRepository) GetViewStatistics() (*models.ViewStatistics, error) {
	stats := &models.ViewStatistics{}

	err := repository.database.QueryRow(QueryViewsToday).Scan(&stats.ViewsToday)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryViewsThisWeek).Scan(&stats.ViewsThisWeek)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryViewsThisMonth).Scan(&stats.ViewsThisMonth)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryViewsThisYear).Scan(&stats.ViewsThisYear)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryTotalViews).Scan(&stats.TotalViews)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (repository *StatisticsRepository) GetLast24HoursViews() (*models.Last24HoursViews, error) {
	rows, err := repository.database.Query(QueryLast24HoursViews)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hours := make([]models.HourlyViewCount, 0)
	for rows.Next() {
		var hour models.HourlyViewCount
		err := rows.Scan(&hour.Hour, &hour.ViewCount)
		if err != nil {
			return nil, err
		}
		hours = append(hours, hour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.Last24HoursViews{Hours: hours}, nil
}

func (repository *StatisticsRepository) GetLastMonthViews() (*models.LastMonthViews, error) {
	rows, err := repository.database.Query(QueryLastMonthViews)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	days := make([]models.DailyViewCount, 0)
	for rows.Next() {
		var day models.DailyViewCount
		err := rows.Scan(&day.Day, &day.ViewCount)
		if err != nil {
			return nil, err
		}
		days = append(days, day)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.LastMonthViews{Days: days}, nil
}

func (repository *StatisticsRepository) GetLastYearViews() (*models.LastYearViews, error) {
	rows, err := repository.database.Query(QueryLastYearViews)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	months := make([]models.MonthlyViewCount, 0)
	for rows.Next() {
		var month models.MonthlyViewCount
		err := rows.Scan(&month.Month, &month.ViewCount)
		if err != nil {
			return nil, err
		}
		months = append(months, month)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.LastYearViews{Months: months}, nil
}

func (repository *StatisticsRepository) GetCategoryViewsDistribution() ([]models.CategoryViewDistribution, error) {
	rows, err := repository.database.Query(QueryCategoryViewsDistribution)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalViews int64
	categories := make([]models.CategoryViewDistribution, 0)

	for rows.Next() {
		var category models.CategoryViewDistribution
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.ViewCount)
		if err != nil {
			return nil, err
		}
		totalViews += category.ViewCount
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range categories {
		if totalViews > 0 {
			categories[i].Percentage = float64(categories[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return categories, nil
}

func (repository *StatisticsRepository) GetMostViewedBlogs(limit int) ([]models.MostViewedBlog, error) {
	rows, err := repository.database.Query(QueryMostViewedBlogs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogs := make([]models.MostViewedBlog, 0)
	for rows.Next() {
		var blog models.MostViewedBlog
		err := rows.Scan(&blog.PostId, &blog.Title, &blog.Slug, &blog.ViewCount, &blog.Author, &blog.CategoryName)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (repository *StatisticsRepository) GetBlogStatistics() (*models.BlogStatistics, error) {
	stats := &models.BlogStatistics{}

	err := repository.database.QueryRow(QueryTotalPublishedBlogs).Scan(&stats.TotalPublishedBlogs)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryTotalDraftedBlogs).Scan(&stats.TotalDraftedBlogs)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryTotalPendingBlogs).Scan(&stats.TotalPendingBlogs)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryTotalReadTime).Scan(&stats.TotalReadTime)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryAverageReadTime).Scan(&stats.AverageReadTime)
	if err != nil {
		return nil, err
	}

	err = repository.database.QueryRow(QueryAverageViews).Scan(&stats.AverageViews)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (repository *StatisticsRepository) GetLongestBlogs(limit int) ([]models.LongestBlog, error) {
	rows, err := repository.database.Query(QueryLongestBlogs, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogs := make([]models.LongestBlog, 0)
	for rows.Next() {
		var blog models.LongestBlog
		err := rows.Scan(&blog.PostId, &blog.Title, &blog.Slug, &blog.ReadTime, &blog.Author, &blog.CategoryName)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (repository *StatisticsRepository) GetCategoryBlogDistribution() ([]models.CategoryBlogDistribution, error) {
	rows, err := repository.database.Query(QueryCategoryBlogDistribution)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalBlogs int64
	categories := make([]models.CategoryBlogDistribution, 0)

	for rows.Next() {
		var category models.CategoryBlogDistribution
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.BlogCount)
		if err != nil {
			return nil, err
		}
		totalBlogs += category.BlogCount
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range categories {
		if totalBlogs > 0 {
			categories[i].Percentage = float64(categories[i].BlogCount) / float64(totalBlogs) * 100
		}
	}

	return categories, nil
}

func (repository *StatisticsRepository) GetCategoryReadTimeDistribution() ([]models.CategoryReadTimeDistribution, error) {
	rows, err := repository.database.Query(QueryCategoryReadTimeDistribution)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalReadTime int64
	categories := make([]models.CategoryReadTimeDistribution, 0)

	for rows.Next() {
		var category models.CategoryReadTimeDistribution
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.TotalReadTime, &category.AverageReadTime)
		if err != nil {
			return nil, err
		}
		totalReadTime += category.TotalReadTime
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range categories {
		if totalReadTime > 0 {
			categories[i].Percentage = float64(categories[i].TotalReadTime) / float64(totalReadTime) * 100
		}
	}

	return categories, nil
}

func (repository *StatisticsRepository) GetCategoryLengthDistribution() ([]models.CategoryLengthDistribution, error) {
	rows, err := repository.database.Query(QueryCategoryLengthDistribution)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalLength int64
	categories := make([]models.CategoryLengthDistribution, 0)

	for rows.Next() {
		var category models.CategoryLengthDistribution
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.TotalLength, &category.AverageLength)
		if err != nil {
			return nil, err
		}
		totalLength += category.TotalLength
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range categories {
		if totalLength > 0 {
			categories[i].Percentage = float64(categories[i].TotalLength) / float64(totalLength) * 100
		}
	}

	return categories, nil
}

func (repository *StatisticsRepository) GetTopUserAgents(limit int) ([]models.UserAgentStat, error) {
	rows, err := repository.database.Query(QueryTopUserAgents, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalViews int64
	userAgents := make([]models.UserAgentStat, 0)

	for rows.Next() {
		var ua models.UserAgentStat
		err := rows.Scan(&ua.UserAgent, &ua.ViewCount)
		if err != nil {
			return nil, err
		}
		totalViews += ua.ViewCount
		userAgents = append(userAgents, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range userAgents {
		if totalViews > 0 {
			userAgents[i].Percentage = float64(userAgents[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return userAgents, nil
}

func (repository *StatisticsRepository) GetDeviceTypeDistribution() ([]models.DeviceTypeStat, error) {
	// Get all user agents first
	rows, err := repository.database.Query(QueryGetAllUserAgents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAgents []string
	for rows.Next() {
		var ua string
		err := rows.Scan(&ua)
		if err != nil {
			return nil, err
		}
		userAgents = append(userAgents, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Parse user agents and count device types
	deviceDistribution := useragent.GetDeviceTypeDistribution(userAgents)

	var totalViews int64
	var deviceStats []models.DeviceTypeStat

	for deviceType, count := range deviceDistribution {
		if count > 0 {
			totalViews += int64(count)
			deviceStats = append(deviceStats, models.DeviceTypeStat{
				DeviceType: deviceType,
				ViewCount:  int64(count),
			})
		}
	}

	// Calculate percentages
	for i := range deviceStats {
		if totalViews > 0 {
			deviceStats[i].Percentage = float64(deviceStats[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return deviceStats, nil
}

func (repository *StatisticsRepository) GetOSDistribution() ([]models.OSStatistic, error) {
	// Get all user agents first
	rows, err := repository.database.Query(QueryGetAllUserAgents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAgents []string
	for rows.Next() {
		var ua string
		err := rows.Scan(&ua)
		if err != nil {
			return nil, err
		}
		userAgents = append(userAgents, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Parse user agents and count operating systems
	osDistribution := useragent.GetOSDistribution(userAgents)

	var totalViews int64
	var osStats []models.OSStatistic

	for os, count := range osDistribution {
		if count > 0 {
			totalViews += int64(count)
			osStats = append(osStats, models.OSStatistic{
				OS:        os,
				ViewCount: int64(count),
			})
		}
	}

	// Calculate percentages
	for i := range osStats {
		if totalViews > 0 {
			osStats[i].Percentage = float64(osStats[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return osStats, nil
}

func (repository *StatisticsRepository) GetBrowserDistribution() ([]models.BrowserStat, error) {
	// Get all user agents first
	rows, err := repository.database.Query(QueryGetAllUserAgents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAgents []string
	for rows.Next() {
		var ua string
		err := rows.Scan(&ua)
		if err != nil {
			return nil, err
		}
		userAgents = append(userAgents, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Parse user agents and count browsers
	browserDistribution := useragent.GetBrowserDistribution(userAgents)

	var totalViews int64
	var browserStats []models.BrowserStat

	for browser, count := range browserDistribution {
		if count > 0 {
			totalViews += int64(count)
			browserStats = append(browserStats, models.BrowserStat{
				Browser:   browser,
				ViewCount: int64(count),
			})
		}
	}

	// Calculate percentages
	for i := range browserStats {
		if totalViews > 0 {
			browserStats[i].Percentage = float64(browserStats[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return browserStats, nil
}

// Author-specific methods
func (repository *StatisticsRepository) GetAuthorStatistics(authorId int64) (*models.AuthorStatistics, error) {
	stats := &models.AuthorStatistics{}

	var totalReadTime int64
	err := repository.database.QueryRow(QueryAuthorStatistics, authorId).Scan(
		&stats.AuthorId, &stats.AuthorName, &stats.TotalBlogs, &stats.TotalViews, &totalReadTime)
	if err != nil {
		return nil, err
	}

	stats.TotalReadTime = totalReadTime
	if stats.TotalBlogs > 0 {
		stats.AverageViews = float64(stats.TotalViews) / float64(stats.TotalBlogs)
		stats.AverageReadTime = float64(stats.TotalReadTime) / float64(stats.TotalBlogs)
	}

	return stats, nil
}

func (repository *StatisticsRepository) GetAuthorViewStatistics(authorId int64) (*models.ViewStatistics, error) {
	stats := &models.ViewStatistics{}

	// Get the total views for posts created by this author
	err := repository.database.QueryRow(QueryAuthorViewsQuery, authorId).Scan(&stats.TotalViews)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (repository *StatisticsRepository) GetAuthorCategoryViewsDistribution(authorId int64) ([]models.CategoryViewDistribution, error) {
	rows, err := repository.database.Query(QueryAuthorCategoryViewsDistribution, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalViews int64
	categories := make([]models.CategoryViewDistribution, 0)

	for rows.Next() {
		var category models.CategoryViewDistribution
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.ViewCount)
		if err != nil {
			return nil, err
		}
		totalViews += category.ViewCount
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range categories {
		if totalViews > 0 {
			categories[i].Percentage = float64(categories[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return categories, nil
}

func (repository *StatisticsRepository) GetAuthorMostViewedBlogs(authorId int64, limit int) ([]models.MostViewedBlog, error) {
	rows, err := repository.database.Query(QueryAuthorMostViewedBlogs, authorId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogs := make([]models.MostViewedBlog, 0)
	for rows.Next() {
		var blog models.MostViewedBlog
		err := rows.Scan(&blog.PostId, &blog.Title, &blog.Slug, &blog.ViewCount, &blog.Author, &blog.CategoryName)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (repository *StatisticsRepository) GetAuthorLongestBlogs(authorId int64, limit int) ([]models.LongestBlog, error) {
	rows, err := repository.database.Query(QueryAuthorLongestBlogs, authorId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogs := make([]models.LongestBlog, 0)
	for rows.Next() {
		var blog models.LongestBlog
		err := rows.Scan(&blog.PostId, &blog.Title, &blog.Slug, &blog.ReadTime, &blog.Author, &blog.CategoryName)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (repository *StatisticsRepository) GetAuthorBlogStatistics(authorId int64) (*models.BlogStatistics, error) {
	stats := &models.BlogStatistics{}

	err := repository.database.QueryRow(QueryAuthorBlogStatistics, authorId, authorId, authorId).Scan(
		&stats.TotalPublishedBlogs, &stats.TotalDraftedBlogs, &stats.TotalPendingBlogs)
	if err != nil {
		return nil, err
	}

	// Get read time and view stats for this author
	err = repository.database.QueryRow(QueryAuthorReadTimeStats, authorId).Scan(&stats.TotalReadTime, &stats.AverageReadTime)
	if err != nil {
		return nil, err
	}

	// Get average views for this author's posts
	err = repository.database.QueryRow(QueryAuthorAverageViews, authorId).Scan(&stats.AverageViews)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (repository *StatisticsRepository) GetAuthorLast24HoursViews(authorId int64) (*models.Last24HoursViews, error) {
	rows, err := repository.database.Query(QueryAuthorLast24HoursViews, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hours := make([]models.HourlyViewCount, 0)
	for rows.Next() {
		var hour models.HourlyViewCount
		err := rows.Scan(&hour.Hour, &hour.ViewCount)
		if err != nil {
			return nil, err
		}
		hours = append(hours, hour)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.Last24HoursViews{Hours: hours}, nil
}

func (repository *StatisticsRepository) GetAuthorLastMonthViews(authorId int64) (*models.LastMonthViews, error) {
	rows, err := repository.database.Query(QueryAuthorLastMonthViews, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	days := make([]models.DailyViewCount, 0)
	for rows.Next() {
		var day models.DailyViewCount
		err := rows.Scan(&day.Day, &day.ViewCount)
		if err != nil {
			return nil, err
		}
		days = append(days, day)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.LastMonthViews{Days: days}, nil
}

func (repository *StatisticsRepository) GetAuthorLastYearViews(authorId int64) (*models.LastYearViews, error) {
	rows, err := repository.database.Query(QueryAuthorLastYearViews, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	months := make([]models.MonthlyViewCount, 0)
	for rows.Next() {
		var month models.MonthlyViewCount
		err := rows.Scan(&month.Month, &month.ViewCount)
		if err != nil {
			return nil, err
		}
		months = append(months, month)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.LastYearViews{Months: months}, nil
}

func (repository *StatisticsRepository) GetAuthorCategoryBlogDistribution(authorId int64) ([]models.CategoryBlogDistribution, error) {
	rows, err := repository.database.Query(QueryAuthorCategoryBlogDistribution, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalBlogs int64
	categories := make([]models.CategoryBlogDistribution, 0)

	for rows.Next() {
		var category models.CategoryBlogDistribution
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.BlogCount)
		if err != nil {
			return nil, err
		}
		totalBlogs += category.BlogCount
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range categories {
		if totalBlogs > 0 {
			categories[i].Percentage = float64(categories[i].BlogCount) / float64(totalBlogs) * 100
		}
	}

	return categories, nil
}

func (repository *StatisticsRepository) GetAuthorCategoryReadTimeDistribution(authorId int64) ([]models.CategoryReadTimeDistribution, error) {
	rows, err := repository.database.Query(QueryAuthorCategoryReadTimeDistribution, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalReadTime int64
	categories := make([]models.CategoryReadTimeDistribution, 0)

	for rows.Next() {
		var category models.CategoryReadTimeDistribution
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.TotalReadTime, &category.AverageReadTime)
		if err != nil {
			return nil, err
		}
		totalReadTime += category.TotalReadTime
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range categories {
		if totalReadTime > 0 {
			categories[i].Percentage = float64(categories[i].TotalReadTime) / float64(totalReadTime) * 100
		}
	}

	return categories, nil
}

func (repository *StatisticsRepository) GetAuthorCategoryLengthDistribution(authorId int64) ([]models.CategoryLengthDistribution, error) {
	rows, err := repository.database.Query(QueryAuthorCategoryLengthDistribution, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalLength int64
	categories := make([]models.CategoryLengthDistribution, 0)

	for rows.Next() {
		var category models.CategoryLengthDistribution
		err := rows.Scan(&category.CategoryId, &category.CategoryName, &category.TotalLength, &category.AverageLength)
		if err != nil {
			return nil, err
		}
		totalLength += category.TotalLength
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range categories {
		if totalLength > 0 {
			categories[i].Percentage = float64(categories[i].TotalLength) / float64(totalLength) * 100
		}
	}

	return categories, nil
}

func (repository *StatisticsRepository) GetAuthorTopUserAgents(authorId int64, limit int) ([]models.UserAgentStat, error) {
	rows, err := repository.database.Query(QueryAuthorTopUserAgents, authorId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalViews int64
	userAgents := make([]models.UserAgentStat, 0)

	for rows.Next() {
		var ua models.UserAgentStat
		err := rows.Scan(&ua.UserAgent, &ua.ViewCount)
		if err != nil {
			return nil, err
		}
		totalViews += ua.ViewCount
		userAgents = append(userAgents, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range userAgents {
		if totalViews > 0 {
			userAgents[i].Percentage = float64(userAgents[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return userAgents, nil
}

func (repository *StatisticsRepository) GetAuthorDeviceTypeDistribution(authorId int64) ([]models.DeviceTypeStat, error) {
	// Get all user agents for this author first
	rows, err := repository.database.Query(QueryAuthorGetAllUserAgents, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAgents []string
	for rows.Next() {
		var ua string
		err := rows.Scan(&ua)
		if err != nil {
			return nil, err
		}
		userAgents = append(userAgents, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Parse user agents and count device types
	deviceDistribution := useragent.GetDeviceTypeDistribution(userAgents)

	var totalViews int64
	var deviceStats []models.DeviceTypeStat

	for deviceType, count := range deviceDistribution {
		if count > 0 {
			totalViews += int64(count)
			deviceStats = append(deviceStats, models.DeviceTypeStat{
				DeviceType: deviceType,
				ViewCount:  int64(count),
			})
		}
	}

	// Calculate percentages
	for i := range deviceStats {
		if totalViews > 0 {
			deviceStats[i].Percentage = float64(deviceStats[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return deviceStats, nil
}

func (repository *StatisticsRepository) GetAuthorOSDistribution(authorId int64) ([]models.OSStatistic, error) {
	// Get all user agents for this author first
	rows, err := repository.database.Query(QueryAuthorGetAllUserAgents, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAgents []string
	for rows.Next() {
		var ua string
		err := rows.Scan(&ua)
		if err != nil {
			return nil, err
		}
		userAgents = append(userAgents, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Parse user agents and count operating systems
	osDistribution := useragent.GetOSDistribution(userAgents)

	var totalViews int64
	var osStats []models.OSStatistic

	for os, count := range osDistribution {
		if count > 0 {
			totalViews += int64(count)
			osStats = append(osStats, models.OSStatistic{
				OS:        os,
				ViewCount: int64(count),
			})
		}
	}

	// Calculate percentages
	for i := range osStats {
		if totalViews > 0 {
			osStats[i].Percentage = float64(osStats[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return osStats, nil
}

func (repository *StatisticsRepository) GetAuthorBrowserDistribution(authorId int64) ([]models.BrowserStat, error) {
	// Get all user agents for this author first
	rows, err := repository.database.Query(QueryAuthorGetAllUserAgents, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userAgents []string
	for rows.Next() {
		var ua string
		err := rows.Scan(&ua)
		if err != nil {
			return nil, err
		}
		userAgents = append(userAgents, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Parse user agents and count browsers
	browserDistribution := useragent.GetBrowserDistribution(userAgents)

	var totalViews int64
	var browserStats []models.BrowserStat

	for browser, count := range browserDistribution {
		if count > 0 {
			totalViews += int64(count)
			browserStats = append(browserStats, models.BrowserStat{
				Browser:   browser,
				ViewCount: int64(count),
			})
		}
	}

	// Calculate percentages
	for i := range browserStats {
		if totalViews > 0 {
			browserStats[i].Percentage = float64(browserStats[i].ViewCount) / float64(totalViews) * 100
		}
	}

	return browserStats, nil
}
