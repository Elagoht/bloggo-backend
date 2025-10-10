package posts

import (
	"bloggo/internal/module/api/posts/models"
)

type PostsAPIService struct {
	repository PostsAPIRepository
}

func NewPostsAPIService(repository PostsAPIRepository) PostsAPIService {
	return PostsAPIService{repository}
}

func (service *PostsAPIService) GetPublishedPosts(page, limit int, categorySlug, tagSlug, authorId, search *string) (*models.APIPostsResponse, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12
	}

	return service.repository.GetPublishedPosts(page, limit, categorySlug, tagSlug, authorId, search)
}

func (service *PostsAPIService) GetPublishedPostBySlug(slug string) (*models.APIPostDetails, error) {
	return service.repository.GetPublishedPostBySlug(slug)
}

func (service *PostsAPIService) TrackView(slug string, userAgent string) error {
	return service.repository.TrackView(slug, userAgent)
}

func (service *PostsAPIService) GetAllViewCounts() (map[string]int64, error) {
	return service.repository.GetAllViewCounts()
}
