package search

import (
	"bloggo/internal/module/search/models"
	"database/sql"
	"fmt"
)

type SearchRepository struct {
	database *sql.DB
}

func NewSearchRepository(database *sql.DB) SearchRepository {
	return SearchRepository{
		database: database,
	}
}

func (repository *SearchRepository) SearchAll(query string, limit int) ([]models.SearchResult, error) {
	results := []models.SearchResult{}
	searchTerm := fmt.Sprintf("%%%s%%", query)
	perTypeLimit := limit / 4 // Distribute limit across 4 types

	if perTypeLimit < 1 {
		perTypeLimit = 1
	}

	// Search tags
	tags, err := repository.searchTags(searchTerm, perTypeLimit)
	if err != nil {
		return nil, err
	}
	results = append(results, tags...)

	// Search categories
	categories, err := repository.searchCategories(searchTerm, perTypeLimit)
	if err != nil {
		return nil, err
	}
	results = append(results, categories...)

	// Search posts
	posts, err := repository.searchPosts(searchTerm, perTypeLimit)
	if err != nil {
		return nil, err
	}
	results = append(results, posts...)

	// Search users
	users, err := repository.searchUsers(searchTerm, perTypeLimit)
	if err != nil {
		return nil, err
	}
	results = append(results, users...)

	return results, nil
}

func (repository *SearchRepository) CountAll(query string) (int, error) {
	searchTerm := fmt.Sprintf("%%%s%%", query)
	total := 0

	// Count tags
	tagCount, err := repository.countTags(searchTerm)
	if err != nil {
		return 0, err
	}
	total += tagCount

	// Count categories
	categoryCount, err := repository.countCategories(searchTerm)
	if err != nil {
		return 0, err
	}
	total += categoryCount

	// Count posts
	postCount, err := repository.countPosts(searchTerm)
	if err != nil {
		return 0, err
	}
	total += postCount

	// Count users
	userCount, err := repository.countUsers(searchTerm)
	if err != nil {
		return 0, err
	}
	total += userCount

	return total, nil
}

func (repository *SearchRepository) searchTags(searchTerm string, limit int) ([]models.SearchResult, error) {
	rows, err := repository.database.Query(QuerySearchTags, searchTerm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanResults(rows)
}

func (repository *SearchRepository) searchCategories(searchTerm string, limit int) ([]models.SearchResult, error) {
	rows, err := repository.database.Query(QuerySearchCategories, searchTerm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanResults(rows)
}

func (repository *SearchRepository) searchPosts(searchTerm string, limit int) ([]models.SearchResult, error) {
	rows, err := repository.database.Query(QuerySearchPosts, searchTerm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanResults(rows)
}

func (repository *SearchRepository) searchUsers(searchTerm string, limit int) ([]models.SearchResult, error) {
	rows, err := repository.database.Query(QuerySearchUsers, searchTerm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return repository.scanResults(rows)
}

func (repository *SearchRepository) countTags(searchTerm string) (int, error) {
	var count int
	err := repository.database.QueryRow(QueryCountSearchTags, searchTerm).Scan(&count)
	return count, err
}

func (repository *SearchRepository) countCategories(searchTerm string) (int, error) {
	var count int
	err := repository.database.QueryRow(QueryCountSearchCategories, searchTerm).Scan(&count)
	return count, err
}

func (repository *SearchRepository) countPosts(searchTerm string) (int, error) {
	var count int
	err := repository.database.QueryRow(QueryCountSearchPosts, searchTerm).Scan(&count)
	return count, err
}

func (repository *SearchRepository) countUsers(searchTerm string) (int, error) {
	var count int
	err := repository.database.QueryRow(QueryCountSearchUsers, searchTerm).Scan(&count)
	return count, err
}

func (repository *SearchRepository) scanResults(rows *sql.Rows) ([]models.SearchResult, error) {
	results := []models.SearchResult{}

	for rows.Next() {
		var result models.SearchResult
		var resultType string

		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Slug,
			&result.AvatarURL,
			&result.CoverURL,
			&resultType,
		)
		if err != nil {
			return nil, err
		}

		result.Type = models.SearchResultType(resultType)
		results = append(results, result)
	}

	return results, rows.Err()
}