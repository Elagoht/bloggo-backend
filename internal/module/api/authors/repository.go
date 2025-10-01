package authors

import (
	"bloggo/internal/module/api/authors/models"
	"bloggo/internal/utils/apierrors"
	"database/sql"
)

type AuthorsAPIRepository struct {
	database *sql.DB
}

func NewAuthorsAPIRepository(database *sql.DB) AuthorsAPIRepository {
	return AuthorsAPIRepository{database}
}

// formatAvatarPath converts database avatar filename to API path format
func formatAvatarPath(filename *string) *string {
	if filename == nil || *filename == "" {
		return nil
	}
	formatted := "/uploads/avatar/" + *filename
	return &formatted
}

func (r *AuthorsAPIRepository) GetAllAuthors() (*models.APIAuthorsResponse, error) {
	rows, err := r.database.Query(QueryAPIGetAllAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	authors := []models.APIAuthorDetails{}
	for rows.Next() {
		var author models.APIAuthorDetails
		err := rows.Scan(
			&author.ID,
			&author.Name,
			&author.Avatar,
			&author.MemberSince,
			&author.PublishedPostCount,
		)
		if err != nil {
			return nil, err
		}
		author.Avatar = formatAvatarPath(author.Avatar)
		authors = append(authors, author)
	}

	return &models.APIAuthorsResponse{
		Authors: authors,
	}, nil
}

func (r *AuthorsAPIRepository) GetAuthorById(id int64) (*models.APIAuthorDetails, error) {
	row := r.database.QueryRow(QueryAPIGetAuthorById, id)

	var author models.APIAuthorDetails
	err := row.Scan(
		&author.ID,
		&author.Name,
		&author.Avatar,
		&author.MemberSince,
		&author.PublishedPostCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierrors.ErrNotFound
		}
		return nil, err
	}

	author.Avatar = formatAvatarPath(author.Avatar)
	return &author, nil
}
