package removal_request

const (
	QueryCreateRemovalRequest = `
	INSERT INTO removal_requests (post_version_id, requested_by, note)
	VALUES (?, ?, ?);`

	QueryGetRemovalRequestList = `
	SELECT
		rr.id, rr.post_version_id, pv.title as post_title,
		u1.id as requested_by_id, u1.name as requested_by_name, u1.avatar as requested_by_avatar,
		rr.note, rr.status,
		u2.id as decided_by_id, u2.name as decided_by_name, u2.avatar as decided_by_avatar,
		rr.decision_note, rr.decided_at, rr.created_at
	FROM removal_requests rr
	JOIN post_versions pv ON rr.post_version_id = pv.id
	JOIN users u1 ON rr.requested_by = u1.id
	LEFT JOIN users u2 ON rr.decided_by = u2.id
	ORDER BY rr.created_at DESC;`

	QueryGetRemovalRequestById = `
	SELECT
		rr.id, rr.post_version_id, pv.title as post_title,
		u_writer.id as writer_id, u_writer.name as writer_name, u_writer.avatar as writer_avatar,
		pv.cover_image as post_cover_url, c.name as post_category,
		u1.id as requested_by_id, u1.name as requested_by_name, u1.avatar as requested_by_avatar,
		rr.note, rr.status,
		u2.id as decided_by_id, u2.name as decided_by_name, u2.avatar as decided_by_avatar,
		rr.decision_note, rr.decided_at, rr.created_at
	FROM removal_requests rr
	JOIN post_versions pv ON rr.post_version_id = pv.id
	JOIN users u_writer ON pv.created_by = u_writer.id
	LEFT JOIN categories c ON pv.category_id = c.id
	JOIN users u1 ON rr.requested_by = u1.id
	LEFT JOIN users u2 ON rr.decided_by = u2.id
	WHERE rr.id = ?;`

	QueryGetUserRemovalRequests = `
	SELECT
		rr.id, rr.post_version_id, pv.title as post_title,
		u1.id as requested_by_id, u1.name as requested_by_name, u1.avatar as requested_by_avatar,
		rr.note, rr.status,
		u2.id as decided_by_id, u2.name as decided_by_name, u2.avatar as decided_by_avatar,
		rr.decision_note, rr.decided_at, rr.created_at
	FROM removal_requests rr
	JOIN post_versions pv ON rr.post_version_id = pv.id
	JOIN users u1 ON rr.requested_by = u1.id
	LEFT JOIN users u2 ON rr.decided_by = u2.id
	WHERE rr.requested_by = ?
	ORDER BY rr.created_at DESC;`

	QueryGetVersionOwnerAndStatus = `
	SELECT created_by, status
	FROM post_versions
	WHERE id = ?
	AND deleted_at IS NULL;`

	QueryCheckExistingRemovalRequest = `
	SELECT COUNT(*)
	FROM removal_requests
	WHERE post_version_id = ?
	AND requested_by = ?
	AND status = 0;`

	QueryApproveRemovalRequest = `
	UPDATE removal_requests
	SET status = 1, decided_by = ?, decision_note = ?, decided_at = CURRENT_TIMESTAMP
	WHERE id = ? AND status = 0;`

	QueryRejectRemovalRequest = `
	UPDATE removal_requests
	SET status = 2, decided_by = ?, decision_note = ?, decided_at = CURRENT_TIMESTAMP
	WHERE id = ? AND status = 0;`

	QueryAutoApproveOtherRemovalRequests = `
	UPDATE removal_requests
	SET status = 1, decided_by = ?, decision_note = ?, decided_at = CURRENT_TIMESTAMP
	WHERE post_version_id = ? AND status = 0 AND id != ?;`

	QueryDeleteRemovalRequestsByVersionId = `
	DELETE FROM removal_requests
	WHERE post_version_id = ?;`

	QueryCheckIfVersionIsCurrentlyPublished = `
	SELECT COUNT(*)
	FROM posts
	WHERE current_version_id = ?
	AND deleted_at IS NULL;`

	QuerySetPostCurrentVersionToNull = `
	UPDATE posts
	SET current_version_id = NULL
	WHERE current_version_id = ?
	AND deleted_at IS NULL;`

	QuerySoftDeleteVersion = `
	UPDATE post_versions
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE id = ?
	AND deleted_at IS NULL;`

	QueryGetVersionCoverImage = `
	SELECT cover_image
	FROM post_versions
	WHERE id = ?
	AND deleted_at IS NULL;`

	QueryCheckImageReferences = `
	SELECT COUNT(*)
	FROM post_versions
	WHERE cover_image = ?
	AND id != ?
	AND deleted_at IS NULL;`

	QueryGetPostIdFromVersionId = `
	SELECT post_id
	FROM post_versions
	WHERE id = ?
	AND deleted_at IS NULL;`

	QueryGetAllVersionsForPost = `
	SELECT cover_image
	FROM post_versions
	WHERE post_id = ?
	AND deleted_at IS NULL;`

	QuerySoftDeleteAllVersionsForPost = `
	UPDATE post_versions
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE post_id = ?
	AND deleted_at IS NULL;`

	QuerySoftDeletePost = `
	UPDATE posts
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE id = ?
	AND deleted_at IS NULL;`

	QueryAutoApproveOtherRemovalRequestsForPost = `
	UPDATE removal_requests
	SET status = 1, decided_by = ?, decision_note = ?, decided_at = CURRENT_TIMESTAMP
	WHERE post_version_id IN (
		SELECT id FROM post_versions WHERE post_id = ?
	)
	AND status = 0
	AND id != ?;`
)