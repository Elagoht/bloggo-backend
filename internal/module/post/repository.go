package post

import (
	"bloggo/internal/module/post/models"
	"database/sql"
)

type PostRepository struct {
	database *sql.DB
}

func NewPostRepository(database *sql.DB) PostRepository {
	return PostRepository{
		database,
	}
}

func (repository *PostRepository) GetPostList() (
	[]models.ResponsePostCard,
	error,
) {
	rows, err := repository.database.Query(QueryPostGetList)
	if err != nil {
		return nil, err
	}

	posts := []models.ResponsePostCard{}
	for rows.Next() {
		var post models.ResponsePostCard

		err := rows.Scan(
			post.PostId,
			post.Author.Name,
			post.Author.Avatar,
			post.Title,
			post.Slug,
			post.CoverImage,
			post.Spot,
			post.Status,
			post.IsActive,
			post.CreatedAt,
			post.UpdatedAt,
			post.Category.Slug,
			post.Category.Id,
			post.Category.Name,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (repository *PostRepository) GetPostBySlug(
	slug string,
) (*models.ResponsePostDetails, error) {
	row := repository.database.QueryRow(QueryPostGetBySlug, slug)

	var post models.ResponsePostDetails
	err := row.Scan(
		post.PostId,
		post.VersionId,
		post.Author.Name,
		post.Author.Email,
		post.Author.Avatar,
		post.Title,
		post.Slug,
		post.Content,
		post.CoverImage,
		post.Description,
		post.Spot,
		post.Status,
		post.StatusChangedAt,
		post.StatusChangedBy,
		post.StatusChangeNote,
		post.IsActive,
		post.CreatedBy,
		post.CreatedAt,
		post.UpdatedAt,
		post.Category.Slug,
		post.Category.Id,
		post.Category.Name,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (repository *PostRepository) CreatePost(
	authorId int64,
) (int64, error) {
	statement, err := repository.database.Prepare(QueryPostCreate)
	if err != nil {
		return 0, err
	}

	result, err := statement.Exec(authorId)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
