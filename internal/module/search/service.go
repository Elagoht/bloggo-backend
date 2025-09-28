package search

import (
	"bloggo/internal/module/search/models"
	"fmt"
	"path/filepath"
	"strings"
)

type SearchService struct {
	repository SearchRepository
}

func NewSearchService(repository SearchRepository) SearchService {
	return SearchService{
		repository: repository,
	}
}

func (service *SearchService) Search(query string, limit int) (models.SearchResponse, error) {
	if query == "" {
		return models.SearchResponse{
			Results: []models.SearchResult{},
			Total:   0,
		}, nil
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	results, err := service.repository.SearchAll(query, limit)
	if err != nil {
		return models.SearchResponse{}, err
	}

	total, err := service.repository.CountAll(query)
	if err != nil {
		return models.SearchResponse{}, err
	}

	// Add path prefixes to URLs
	for i := range results {
		// Add avatar URL prefix if avatar exists
		if results[i].AvatarURL != nil && *results[i].AvatarURL != "" {
			avatarPath := fmt.Sprintf("/uploads/avatar/%s", *results[i].AvatarURL)
			results[i].AvatarURL = &avatarPath
		}
		// Format cover URL if cover exists
		if results[i].CoverURL != nil && *results[i].CoverURL != "" {
			nameWithoutExt := strings.TrimSuffix(*results[i].CoverURL, filepath.Ext(*results[i].CoverURL))
			coverPath := fmt.Sprintf("/uploads/cover/%s", nameWithoutExt)
			results[i].CoverURL = &coverPath
		}
	}

	return models.SearchResponse{
		Results: results,
		Total:   total,
	}, nil
}