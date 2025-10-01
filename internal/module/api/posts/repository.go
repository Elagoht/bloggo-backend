package posts

import (
	"bloggo/internal/module/api/posts/models"
	"bloggo/internal/utils/apierrors"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

type PostsAPIRepository struct {
	database *sql.DB
}

func NewPostsAPIRepository(database *sql.DB) PostsAPIRepository {
	return PostsAPIRepository{database}
}

// formatCoverImagePath converts database cover image filename to API path format
func formatCoverImagePath(filename *string) *string {
	if filename == nil || *filename == "" {
		return nil
	}
	nameWithoutExt := strings.TrimSuffix(*filename, filepath.Ext(*filename))
	formatted := "/uploads/cover/" + nameWithoutExt
	return &formatted
}

// formatAvatarPath converts database avatar filename to API path format
func formatAvatarPath(filename *string) *string {
	if filename == nil || *filename == "" {
		return nil
	}
	formatted := "/uploads/avatar/" + *filename
	return &formatted
}

func (r *PostsAPIRepository) GetPublishedPosts(page, limit int, categorySlug, tagSlug, authorId, search *string) (*models.APIPostsResponse, error) {
	var whereClauses []string
	var args []any

	// Add category filter
	if categorySlug != nil && *categorySlug != "" {
		whereClauses = append(whereClauses, "c.slug = ?")
		args = append(args, *categorySlug)
	}

	// Add tag filter
	if tagSlug != nil && *tagSlug != "" {
		whereClauses = append(whereClauses, "EXISTS (SELECT 1 FROM post_tags pt JOIN tags t ON t.id = pt.tag_id WHERE pt.post_id = p.id AND t.slug = ? AND t.deleted_at IS NULL)")
		args = append(args, *tagSlug)
	}

	// Add author filter
	if authorId != nil && *authorId != "" {
		whereClauses = append(whereClauses, "u.id = ?")
		args = append(args, *authorId)
	}

	// Add search filter
	if search != nil && *search != "" {
		searchTerm := "%" + *search + "%"
		whereClauses = append(whereClauses, "(pv.title LIKE ? OR pv.content LIKE ? OR pv.spot LIKE ?)")
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	// Build WHERE clause
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "AND " + strings.Join(whereClauses, " AND ")
	}

	// Build ORDER BY and pagination
	orderClause := "ORDER BY pv.updated_at DESC LIMIT ? OFFSET ?"
	offset := (page - 1) * limit

	// Count total
	countQuery := fmt.Sprintf(QueryAPICountPublishedPosts, whereClause)
	var total int64
	err := r.database.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get posts
	postsQuery := fmt.Sprintf(QueryAPIGetPublishedPosts, whereClause, orderClause)
	queryArgs := append(args, limit, offset)

	rows, err := r.database.Query(postsQuery, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.APIPostCard{}
	for rows.Next() {
		var post models.APIPostCard
		var rawCoverImage *string

		err := rows.Scan(
			&post.Slug,
			&post.Title,
			&post.Description,
			&post.Spot,
			&rawCoverImage,
			&post.ReadCount,
			&post.ReadTime,
			&post.PublishedAt,
			&post.Author.ID,
			&post.Author.Name,
			&post.Author.Avatar,
			&post.Category.Slug,
			&post.Category.Name,
		)
		if err != nil {
			return nil, err
		}

		post.CoverImage = formatCoverImagePath(rawCoverImage)
		post.Author.Avatar = formatAvatarPath(post.Author.Avatar)

		posts = append(posts, post)
	}

	// For each post, get its tags
	for i := range posts {
		tags, err := r.GetPostTagsBySlug(posts[i].Slug)
		if err != nil {
			return nil, err
		}
		posts[i].Tags = tags
	}

	return &models.APIPostsResponse{
		Data:  posts,
		Page:  page,
		Take:  limit,
		Total: total,
	}, nil
}

func (r *PostsAPIRepository) GetPublishedPostBySlug(slug string) (*models.APIPostDetails, error) {
	row := r.database.QueryRow(QueryAPIGetPublishedPostBySlug, slug)

	var post models.APIPostDetails
	var rawCoverImage *string
	var postId int64
	var categoryDescription *string

	err := row.Scan(
		&post.Slug,
		&post.Title,
		&post.Content,
		&post.Description,
		&post.Spot,
		&rawCoverImage,
		&post.ReadCount,
		&post.ReadTime,
		&post.PublishedAt,
		&post.UpdatedAt,
		&postId,
		&post.Author.ID,
		&post.Author.Name,
		&post.Author.Avatar,
		&post.Category.Slug,
		&post.Category.Name,
		&categoryDescription,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierrors.ErrNotFound
		}
		return nil, err
	}

	post.CoverImage = formatCoverImagePath(rawCoverImage)
	post.Author.Avatar = formatAvatarPath(post.Author.Avatar)

	// Get tags
	tags, err := r.GetPostTags(postId)
	if err != nil {
		return nil, err
	}
	post.Tags = tags

	return &post, nil
}

func (r *PostsAPIRepository) GetPostTags(postId int64) ([]models.APITag, error) {
	rows, err := r.database.Query(QueryAPIGetPostTags, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.APITag{}
	for rows.Next() {
		var tag models.APITag
		err := rows.Scan(&tag.Slug, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *PostsAPIRepository) GetPostTagsBySlug(slug string) ([]models.APITag, error) {
	// First get post ID
	var postId int64
	row := r.database.QueryRow(`
		SELECT p.id FROM posts p
		JOIN post_versions pv ON pv.id = p.current_version_id
		WHERE pv.slug = ? AND p.deleted_at IS NULL AND pv.deleted_at IS NULL
		LIMIT 1`, slug)

	err := row.Scan(&postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.APITag{}, nil
		}
		return nil, err
	}

	return r.GetPostTags(postId)
}

func (r *PostsAPIRepository) TrackView(slug string, userAgent string) error {
	// Get post ID from slug
	var postId int64
	row := r.database.QueryRow(`
		SELECT p.id FROM posts p
		JOIN post_versions pv ON pv.id = p.current_version_id
		WHERE pv.slug = ? AND p.deleted_at IS NULL AND pv.deleted_at IS NULL AND pv.status = 5
		LIMIT 1`, slug)

	err := row.Scan(&postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierrors.ErrNotFound
		}
		return err
	}

	// Track view
	transaction, err := r.database.Begin()
	if err != nil {
		return err
	}

	// Insert view record
	_, err = transaction.Exec(`INSERT INTO post_views (post_id, user_agent) VALUES (?, ?)`, postId, userAgent)
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Increment read count
	_, err = transaction.Exec(`UPDATE posts SET read_count = read_count + 1 WHERE id = ? AND deleted_at IS NULL`, postId)
	if err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit()
}
