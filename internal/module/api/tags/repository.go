package tags

import (
	"bloggo/internal/module/api/tags/models"
	"bloggo/internal/utils/apierrors"
	"database/sql"
)

type TagsAPIRepository struct {
	database *sql.DB
}

func NewTagsAPIRepository(database *sql.DB) TagsAPIRepository {
	return TagsAPIRepository{database}
}

func (r *TagsAPIRepository) GetAllTags() (*models.APITagsResponse, error) {
	rows, err := r.database.Query(QueryAPIGetAllTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.APITagDetails{}
	for rows.Next() {
		var tag models.APITagDetails
		err := rows.Scan(&tag.Slug, &tag.Name, &tag.PostCount)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return &models.APITagsResponse{
		Tags: tags,
	}, nil
}

func (r *TagsAPIRepository) GetTagBySlug(slug string) (*models.APITagDetails, error) {
	row := r.database.QueryRow(QueryAPIGetTagBySlug, slug)

	var tag models.APITagDetails
	err := row.Scan(&tag.Slug, &tag.Name, &tag.PostCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierrors.ErrNotFound
		}
		return nil, err
	}

	return &tag, nil
}
