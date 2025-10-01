package tags

import (
	"bloggo/internal/module/api/tags/models"
)

type TagsAPIService struct {
	repository TagsAPIRepository
}

func NewTagsAPIService(repository TagsAPIRepository) TagsAPIService {
	return TagsAPIService{repository}
}

func (service *TagsAPIService) GetAllTags() (*models.APITagsResponse, error) {
	return service.repository.GetAllTags()
}

func (service *TagsAPIService) GetTagBySlug(slug string) (*models.APITagDetails, error) {
	return service.repository.GetTagBySlug(slug)
}
