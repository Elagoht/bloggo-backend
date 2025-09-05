package statistics

const (
	// View Statistics Queries
	QueryViewsToday = `
	SELECT COUNT(*)
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE DATE(pv.viewed_at) = DATE('now')
	AND p.deleted_at IS NULL`

	QueryViewsThisWeek = `
	SELECT COUNT(*)
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE DATE(pv.viewed_at) >= DATE('now', 'weekday 0', '-6 days')
	AND p.deleted_at IS NULL`

	QueryViewsThisMonth = `
	SELECT COUNT(*)
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE DATE(pv.viewed_at) >= DATE('now', 'start of month')
	AND p.deleted_at IS NULL`

	QueryViewsThisYear = `
	SELECT COUNT(*)
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE DATE(pv.viewed_at) >= DATE('now', 'start of year')
	AND p.deleted_at IS NULL`

	QueryTotalViews = `
	SELECT COUNT(*)
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE p.deleted_at IS NULL`

	// Last 24 Hours Views by Hour
	QueryLast24HoursViews = `
	SELECT
	CAST(strftime('%H', pv.viewed_at) AS INTEGER) as hour,
		COUNT(*) as view_count
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE pv.viewed_at >= datetime('now', '-24 hours')
	AND p.deleted_at IS NULL
	GROUP BY hour
	ORDER BY hour`
	// Last Month Views by Day
	QueryLastMonthViews = `
	SELECT
		CAST(strftime('%d', pv.viewed_at) AS INTEGER) as day,
		COUNT(*) as view_count
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE DATE(pv.viewed_at) >= DATE('now', 'start of month', '-1 month')
	AND DATE(pv.viewed_at) < DATE('now', 'start of month')
	AND p.deleted_at IS NULL
	GROUP BY day
	ORDER BY day`
	// Last 12 Months Views by Month
	QueryLastYearViews = `
	SELECT
		CAST(strftime('%m', pv.viewed_at) AS INTEGER) as month,
		COUNT(*) as view_count
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE pv.viewed_at >= datetime('now', '-12 months')
	AND p.deleted_at IS NULL
	GROUP BY month
	ORDER BY month`
	// Category Views Distribution
	QueryCategoryViewsDistribution = `
	SELECT
	c.id as category_id,
		c.name as category_name,
		COUNT(pv.id) as view_count
	FROM categories c
	LEFT JOIN post_versions ver ON ver.category_id = c.id
	LEFT JOIN posts p ON p.current_version_id = ver.id
	LEFT JOIN post_views pv ON pv.post_id = p.id
	WHERE c.deleted_at IS NULL
	AND (p.deleted_at IS NULL OR p.deleted_at IS NULL)
	AND (ver.deleted_at IS NULL OR ver.deleted_at IS NULL)
	AND ver.status = 5
	GROUP BY c.id, c.name
	ORDER BY view_count DESC`

	// Most Viewed Blogs
	QueryMostViewedBlogs = `
	SELECT
	p.id as post_id,
		ver.title,
		ver.slug,
		COUNT(pv.id) as view_count,
		u.name as author,
		c.name as category_name
	FROM posts p
	JOIN post_versions ver ON p.current_version_id = ver.id
	JOIN users u ON ver.created_by = u.id
	LEFT JOIN categories c ON ver.category_id = c.id
	LEFT JOIN post_views pv ON pv.post_id = p.id
	WHERE p.deleted_at IS NULL
	AND ver.deleted_at IS NULL
	AND ver.status = 5
	GROUP BY p.id, ver.title, ver.slug, u.name, c.name
	ORDER BY view_count DESC
	LIMIT ?`

	// Blog Statistics
	QueryTotalPublishedBlogs = `
	SELECT COUNT(DISTINCT p.id)
	FROM posts p
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE p.deleted_at IS NULL
	AND ver.deleted_at IS NULL
	AND ver.status = 5`

	QueryTotalDraftedBlogs = `
	SELECT COUNT(DISTINCT ver.id)
	FROM post_versions ver
	JOIN posts p ON ver.post_id = p.id
	WHERE ver.deleted_at IS NULL
	AND p.deleted_at IS NULL
	AND ver.status = 0`

	QueryTotalPendingBlogs = `
	SELECT COUNT(DISTINCT ver.id)
	FROM post_versions ver
	JOIN posts p ON ver.post_id = p.id
	WHERE ver.deleted_at IS NULL
	AND p.deleted_at IS NULL
	AND ver.status = 1`

	QueryTotalReadTime = `
	SELECT COALESCE(SUM(ver.read_time), 0)
	FROM posts p
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE p.deleted_at IS NULL
	AND ver.deleted_at IS NULL
	AND ver.status = 5
	AND ver.read_time IS NOT NULL`

	QueryAverageReadTime = `
	SELECT COALESCE(AVG(CAST(ver.read_time AS REAL)), 0)
	FROM posts p
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE p.deleted_at IS NULL
	AND ver.deleted_at IS NULL
	AND ver.status = 5
	AND ver.read_time IS NOT NULL`

	QueryAverageViews = `
	SELECT COALESCE(AVG(CAST(view_counts.view_count AS REAL)), 0)
	FROM (
		SELECT p.id, COUNT(pv.id) as view_count
		FROM posts p
		JOIN post_versions ver ON p.current_version_id = ver.id
		LEFT JOIN post_views pv ON pv.post_id = p.id
		WHERE p.deleted_at IS NULL
		AND ver.deleted_at IS NULL
		AND ver.status = 5
		GROUP BY p.id
	) as view_counts`

	// Longest Blogs
	QueryLongestBlogs = `
	SELECT
	p.id as post_id,
		ver.title,
		ver.slug,
		ver.read_time,
		u.name as author,
		c.name as category_name
	FROM posts p
	JOIN post_versions ver ON p.current_version_id = ver.id
	JOIN users u ON ver.created_by = u.id
	LEFT JOIN categories c ON ver.category_id = c.id
	WHERE p.deleted_at IS NULL
	AND ver.deleted_at IS NULL
	AND ver.status = 5
	AND ver.read_time IS NOT NULL
	ORDER BY ver.read_time DESC
	LIMIT ?`

	// Category Blog Distribution
	QueryCategoryBlogDistribution = `
	SELECT
	c.id as category_id,
		c.name as category_name,
		COUNT(DISTINCT p.id) as blog_count
	FROM categories c
	LEFT JOIN post_versions ver ON ver.category_id = c.id
	LEFT JOIN posts p ON p.current_version_id = ver.id
	WHERE c.deleted_at IS NULL
	AND (p.deleted_at IS NULL OR p.deleted_at IS NULL)
	AND (ver.deleted_at IS NULL OR ver.deleted_at IS NULL)
	AND ver.status = 5
	GROUP BY c.id, c.name
	ORDER BY blog_count DESC`

	// Category Read Time Distribution
	QueryCategoryReadTimeDistribution = `
	SELECT
	c.id as category_id,
		c.name as category_name,
		COALESCE(SUM(ver.read_time), 0) as total_read_time,
		COALESCE(AVG(CAST(ver.read_time AS REAL)), 0) as average_read_time
	FROM categories c
	LEFT JOIN post_versions ver ON ver.category_id = c.id
	LEFT JOIN posts p ON p.current_version_id = ver.id
	WHERE c.deleted_at IS NULL
	AND (p.deleted_at IS NULL OR p.deleted_at IS NULL)
	AND (ver.deleted_at IS NULL OR ver.deleted_at IS NULL)
	AND ver.status = 5
	AND ver.read_time IS NOT NULL
	GROUP BY c.id, c.name
	ORDER BY total_read_time DESC`

	// Category Length Distribution
	QueryCategoryLengthDistribution = `
	SELECT
	c.id as category_id,
		c.name as category_name,
		COALESCE(SUM(LENGTH(ver.content)), 0) as total_length,
		COALESCE(AVG(CAST(LENGTH(ver.content) AS REAL)), 0) as average_length
	FROM categories c
	LEFT JOIN post_versions ver ON ver.category_id = c.id
	LEFT JOIN posts p ON p.current_version_id = ver.id
	WHERE c.deleted_at IS NULL
	AND (p.deleted_at IS NULL OR p.deleted_at IS NULL)
	AND (ver.deleted_at IS NULL OR ver.deleted_at IS NULL)
	AND ver.status = 5
	AND ver.content IS NOT NULL
	GROUP BY c.id, c.name
	ORDER BY total_length DESC`

	// User Agent Statistics
	QueryTopUserAgents = `
	SELECT
	pv.user_agent,
		COUNT(*) as view_count
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE pv.user_agent IS NOT NULL
	AND pv.user_agent != ''
	AND p.deleted_at IS NULL
	GROUP BY pv.user_agent
	ORDER BY view_count DESC
	LIMIT ?`

	// Author-specific queries (add WHERE clause for author filtering)
	QueryAuthorViewsToday     = QueryViewsToday + ` AND ver.created_by = ?`
	QueryAuthorViewsThisWeek  = QueryViewsThisWeek + ` AND ver.created_by = ?`
	QueryAuthorViewsThisMonth = QueryViewsThisMonth + ` AND ver.created_by = ?`
	QueryAuthorViewsThisYear  = QueryViewsThisYear + ` AND ver.created_by = ?`
	QueryAuthorTotalViews     = QueryTotalViews + ` AND ver.created_by = ?`

	QueryAuthorStatistics = `
	SELECT
		u.id as author_id,
		u.name as author_name,
		COUNT(DISTINCT CASE WHEN ver.status = 5 AND p.deleted_at IS NULL AND ver.deleted_at IS NULL THEN p.id END) as total_blogs,
		COUNT(CASE WHEN ver.status = 5 AND p.deleted_at IS NULL AND ver.deleted_at IS NULL THEN pv.id END) as total_views,
		COALESCE(SUM(CASE WHEN ver.status = 5 AND p.deleted_at IS NULL AND ver.deleted_at IS NULL THEN ver.read_time END), 0) as total_read_time
	FROM users u
	LEFT JOIN post_versions ver
		ON ver.created_by = u.id
	LEFT JOIN posts p
		ON p.current_version_id = ver.id
	LEFT JOIN post_views pv
		ON pv.post_id = p.id
	WHERE u.id = ?
	GROUP BY u.id, u.name`

	QueryAuthorCategoryViewsDistribution = `
	SELECT
		c.id as category_id,
		c.name as category_name,
		COUNT(pv.id) as view_count
	FROM categories c
	LEFT JOIN post_versions ver
		ON ver.category_id = c.id
		AND ver.created_by = ?
	LEFT JOIN posts p
		ON p.current_version_id = ver.id
	LEFT JOIN post_views pv
		ON pv.post_id = p.id
	WHERE c.deleted_at IS NULL
		AND (
			p.deleted_at IS NULL
			OR p.deleted_at IS NULL
		)
		AND (
			ver.deleted_at IS NULL
			OR ver.deleted_at IS NULL
		)
		AND ver.status = 5
	GROUP BY c.id, c.name
	ORDER BY view_count DESC`

	QueryAuthorMostViewedBlogs = `
	SELECT
	p.id as post_id,
		ver.title,
		ver.slug,
		COUNT(pv.id) as view_count,
		u.name as author,
		c.name as category_name
	FROM posts p
	JOIN post_versions ver
		ON p.current_version_id = ver.id
		AND ver.created_by = ?
	JOIN users u
		ON ver.created_by = u.id
	LEFT JOIN categories c
		ON ver.category_id = c.id
	LEFT JOIN post_views pv
		ON pv.post_id = p.id
	WHERE p.deleted_at IS NULL
	AND ver.deleted_at IS NULL
	AND ver.status = 5
	GROUP BY p.id, ver.title, ver.slug, u.name, c.name
	ORDER BY view_count DESC
	LIMIT ?`

	QueryAuthorLongestBlogs = `
	SELECT
	p.id as post_id,
		ver.title,
		ver.slug,
		ver.read_time,
		u.name as author,
		c.name as category_name
	FROM posts p
	JOIN post_versions ver
		ON p.current_version_id = ver.id
		AND ver.created_by = ?
	JOIN users u
		ON ver.created_by = u.id
	LEFT JOIN categories c
		ON ver.category_id = c.id
	WHERE p.deleted_at IS NULL
		AND ver.deleted_at IS NULL
		AND ver.status = 5
		AND ver.read_time IS NOT NULL
	ORDER BY ver.read_time DESC
	LIMIT ?`

	QueryAuthorBlogStatistics = `
	SELECT
	(SELECT COUNT(DISTINCT p.id) FROM posts p
		JOIN post_versions ver
		ON p.current_version_id = ver.id
		WHERE ver.created_by = ?
			AND p.deleted_at IS NULL
			AND ver.deleted_at IS NULL
			AND ver.status = 5) as published_blogs,
		(SELECT COUNT(DISTINCT ver.id) FROM post_versions ver
		JOIN posts p
		ON ver.post_id = p.id
		WHERE ver.created_by = ?
			AND ver.deleted_at IS NULL
			AND p.deleted_at IS NULL
			AND ver.status = 0) as drafted_blogs,
		(SELECT COUNT(DISTINCT ver.id) FROM post_versions ver
		JOIN posts p
		ON ver.post_id = p.id
		WHERE ver.created_by = ?
			AND ver.deleted_at IS NULL
			AND p.deleted_at IS NULL
			AND ver.status = 1) as pending_blogs`

	// User Agent Analysis Queries
	QueryGetAllUserAgents = `
	SELECT pv.user_agent
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	WHERE pv.user_agent IS NOT NULL 
	AND pv.user_agent != ''
	AND p.deleted_at IS NULL`

	// Author-specific queries for repository
	QueryAuthorViewsQuery = `
	SELECT COUNT(*) 
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE p.deleted_at IS NULL AND ver.created_by = ?`

	QueryAuthorReadTimeStats = `
	SELECT COALESCE(SUM(ver.read_time), 0), COALESCE(AVG(CAST(ver.read_time AS REAL)), 0)
	FROM posts p
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE p.deleted_at IS NULL 
	AND ver.deleted_at IS NULL
	AND ver.status = 5
	AND ver.created_by = ?
	AND ver.read_time IS NOT NULL`

	QueryAuthorAverageViews = `
	SELECT COALESCE(AVG(CAST(view_counts.view_count AS REAL)), 0)
	FROM (
		SELECT p.id, COUNT(pv.id) as view_count
		FROM posts p
		JOIN post_versions ver ON p.current_version_id = ver.id
		LEFT JOIN post_views pv ON pv.post_id = p.id
		WHERE p.deleted_at IS NULL 
		AND ver.deleted_at IS NULL
		AND ver.status = 5
		AND ver.created_by = ?
		GROUP BY p.id
	) as view_counts`

	// Author-specific Last 24 Hours Views
	QueryAuthorLast24HoursViews = `
	SELECT
		CAST(strftime('%H', pv.viewed_at) AS INTEGER) as hour,
		COUNT(*) as view_count
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE pv.viewed_at >= datetime('now', '-24 hours')
		AND p.deleted_at IS NULL
		AND ver.created_by = ?
	GROUP BY hour
	ORDER BY hour`

	// Author-specific Last Month Views by Day
	QueryAuthorLastMonthViews = `
	SELECT
		CAST(strftime('%d', pv.viewed_at) AS INTEGER) as day,
		COUNT(*) as view_count
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE DATE(pv.viewed_at) >= DATE('now', 'start of month', '-1 month')
		AND DATE(pv.viewed_at) < DATE('now', 'start of month')
		AND p.deleted_at IS NULL
		AND ver.created_by = ?
	GROUP BY day
	ORDER BY day`

	// Author-specific Last 12 Months Views by Month
	QueryAuthorLastYearViews = `
	SELECT
		CAST(strftime('%m', pv.viewed_at) AS INTEGER) as month,
		COUNT(*) as view_count
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE pv.viewed_at >= datetime('now', '-12 months')
		AND p.deleted_at IS NULL
		AND ver.created_by = ?
	GROUP BY month
	ORDER BY month`

	// Author-specific Category Blog Distribution
	QueryAuthorCategoryBlogDistribution = `
	SELECT
		c.id as category_id,
		c.name as category_name,
		COUNT(DISTINCT p.id) as blog_count
	FROM categories c
	LEFT JOIN post_versions ver 
		ON ver.category_id = c.id
		AND ver.created_by = ?
	LEFT JOIN posts p 
		ON p.current_version_id = ver.id
	WHERE c.deleted_at IS NULL
		AND (p.deleted_at IS NULL OR p.deleted_at IS NULL)
		AND (ver.deleted_at IS NULL OR ver.deleted_at IS NULL)
		AND ver.status = 5
	GROUP BY c.id, c.name
	ORDER BY blog_count DESC`

	// Author-specific Category Read Time Distribution
	QueryAuthorCategoryReadTimeDistribution = `
	SELECT
		c.id as category_id,
		c.name as category_name,
		COALESCE(SUM(ver.read_time), 0) as total_read_time,
		COALESCE(AVG(CAST(ver.read_time AS REAL)), 0) as average_read_time
	FROM categories c
	LEFT JOIN post_versions ver 
		ON ver.category_id = c.id
		AND ver.created_by = ?
	LEFT JOIN posts p 
		ON p.current_version_id = ver.id
	WHERE c.deleted_at IS NULL
		AND (p.deleted_at IS NULL OR p.deleted_at IS NULL)
		AND (ver.deleted_at IS NULL OR ver.deleted_at IS NULL)
		AND ver.status = 5
		AND ver.read_time IS NOT NULL
	GROUP BY c.id, c.name
	ORDER BY total_read_time DESC`

	// Author-specific Category Length Distribution
	QueryAuthorCategoryLengthDistribution = `
	SELECT
		c.id as category_id,
		c.name as category_name,
		COALESCE(SUM(LENGTH(ver.content)), 0) as total_length,
		COALESCE(AVG(CAST(LENGTH(ver.content) AS REAL)), 0) as average_length
	FROM categories c
	LEFT JOIN post_versions ver 
		ON ver.category_id = c.id
		AND ver.created_by = ?
	LEFT JOIN posts p 
		ON p.current_version_id = ver.id
	WHERE c.deleted_at IS NULL
		AND (p.deleted_at IS NULL OR p.deleted_at IS NULL)
		AND (ver.deleted_at IS NULL OR ver.deleted_at IS NULL)
		AND ver.status = 5
		AND ver.content IS NOT NULL
	GROUP BY c.id, c.name
	ORDER BY total_length DESC`

	// Author-specific User Agent Queries
	QueryAuthorTopUserAgents = `
	SELECT
		pv.user_agent,
		COUNT(*) as view_count
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE pv.user_agent IS NOT NULL
		AND pv.user_agent != ''
		AND p.deleted_at IS NULL
		AND ver.created_by = ?
	GROUP BY pv.user_agent
	ORDER BY view_count DESC
	LIMIT ?`

	QueryAuthorGetAllUserAgents = `
	SELECT pv.user_agent
	FROM post_views pv
	JOIN posts p ON pv.post_id = p.id
	JOIN post_versions ver ON p.current_version_id = ver.id
	WHERE pv.user_agent IS NOT NULL 
		AND pv.user_agent != ''
		AND p.deleted_at IS NULL
		AND ver.created_by = ?`
)
