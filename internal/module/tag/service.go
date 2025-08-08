package tag

import (
	"bloggo/internal/module/tag/models"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
)

type TagService struct {
	repository TagRepository
}

func NewTagService(repository TagRepository) TagService {
	return TagService{
		repository,
	}
}

func (service *TagService) TagCreate(
	model *models.RequestTagCreate,
) (*models.ResponseTagCreated, error) {
	id, err := service.repository.TagCreate(
		models.ToCreateTagParams(model),
	)
	if err != nil {
		return nil, err
	}

	return &models.ResponseTagCreated{
		Id: id,
	}, nil
}

func (service *TagService) GetTagBySlug(
	slug string,
) (*models.ResponseTagDetails, error) {
	return service.repository.GetTagBySlug(slug)
}

func (service *TagService) GetCategories(
	pagination *pagination.PaginationOptions,
	search *filter.SearchOptions,
) ([]models.ResponseTagCard, error) {
	return service.repository.GetCategories(pagination, search)
}

func (service *TagService) TagUpdate(
	slug string,
	model *models.RequestTagUpdate,
) error {
	return service.repository.TagUpdate(
		slug,
		models.ToUpdateTagParams(model),
	)
}

func (service *TagService) TagDelete(
	slug string,
) error {
	return service.repository.TagDelete(slug)
}
