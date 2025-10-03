package dashboard

const (
	// Get pending versions (status = 1 is pending/submitted for review)
	QueryGetPendingVersions = `
	SELECT
		pv.id, pv.post_id, pv.title, pv.created_by, u.name as author_name, u.avatar as author_avatar, pv.created_at
	FROM post_versions pv
	JOIN users u ON pv.created_by = u.id
	WHERE pv.status = 1
		AND pv.deleted_at IS NULL
	ORDER BY pv.created_at DESC
	LIMIT 10`
	// Get recent published activity (status = 5 is published)
	QueryGetRecentActivity = `
	SELECT
		pv.id, pv.title, pv.status_changed_at as published_at
	FROM post_versions pv
	WHERE pv.status = 5
		AND pv.deleted_at IS NULL
	ORDER BY pv.status_changed_at DESC
	LIMIT 10`
	// Get publishing rate for this week
	QueryGetPublishingRateWeek = `
	SELECT COUNT(*)
	FROM post_versions pv
	WHERE pv.status = 5
		AND pv.status_changed_at >= datetime('now', '-7 days')
		AND pv.deleted_at IS NULL`
	// Get publishing rate for this month
	QueryGetPublishingRateMonth = `
	SELECT COUNT(*)
	FROM post_versions pv
	WHERE pv.status = 5
		AND pv.status_changed_at >= datetime('now', '-30 days')
		AND pv.deleted_at IS NULL`
	// Get publishing rate for today
	QueryGetPublishingRateDay = `
	SELECT COUNT(*)
	FROM post_versions pv
	WHERE pv.status = 5
		AND pv.status_changed_at >= datetime('now', 'start of day')
		AND pv.deleted_at IS NULL`
	// Get publishing rate for this year
	QueryGetPublishingRateYear = `
	SELECT COUNT(*)
	FROM post_versions pv
	WHERE pv.status = 5
		AND pv.status_changed_at >= datetime('now', 'start of year')
		AND pv.deleted_at IS NULL`
	// Get author performance (published posts)
	QueryGetAuthorPerformance = `
	SELECT
		u.id as author_id, u.name as author_name, u.avatar as author_avatar, COUNT(*) as post_count
	FROM post_versions pv
	JOIN users u ON pv.created_by = u.id
	WHERE pv.status = 5
		AND pv.deleted_at IS NULL
	GROUP BY u.id, u.name
	ORDER BY post_count DESC
	LIMIT 3`
	// Get total draft count
	QueryGetTotalDraftCount = `
	SELECT COUNT(*)
	FROM post_versions pv
	WHERE pv.status = 0
		AND pv.deleted_at IS NULL`
	// Get drafts by author
	QueryGetDraftsByAuthor = `
	SELECT
		u.id as author_id, u.name as author_name, u.avatar as author_avatar, COUNT(*) as draft_count
	FROM post_versions pv
	JOIN users u ON pv.created_by = u.id
	WHERE pv.status = 0
		AND pv.deleted_at IS NULL
	GROUP BY u.id, u.name, u.avatar
	ORDER BY draft_count DESC
	LIMIT 5`
	// Get popular tags
	QueryGetPopularTags = `
	SELECT
		t.id, t.name, t.slug, COUNT(*) as usage
	FROM tags t
	JOIN post_tags pt ON t.id = pt.tag_id
	JOIN posts p ON pt.post_id = p.id
	WHERE t.deleted_at IS NULL
		AND p.deleted_at IS NULL
	GROUP BY t.id, t.name
	ORDER BY usage DESC
	LIMIT 10`
)
