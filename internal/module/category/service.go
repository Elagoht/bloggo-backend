package category

import (
	"bloggo/internal/module/category/models"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
)

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repository CategoryRepository) CategoryService {
	return CategoryService{
		repository,
	}
}

func (service *CategoryService) CategoryCreate(
	model *models.RequestCategoryCreate,
) (*models.ResponseCategoryCreated, error) {
	id, err := service.repository.CategoryCreate(
		models.ToCreateCategoryParams(model),
	)
	if err != nil {
		return nil, err
	}

	return &models.ResponseCategoryCreated{
		Id: id,
	}, nil
}

func (service *CategoryService) GetCategoryBySlug(
	slug string,
) (*models.ResponseCategoryDetails, error) {
	return service.repository.GetCategoryBySlug(slug)
}

func (service *CategoryService) GetCategories(
	pagination *pagination.PaginationOptions,
	search *filter.SearchOptions,
) ([]models.ResponseCategoryCard, error) {
	return service.repository.GetCategories(pagination, search)
}

func (service *CategoryService) CategoryUpdate(
	slug string,
	model *models.RequestCategoryUpdate,
) error {
	return service.repository.CategoryUpdate(
		slug,
		models.ToUpdateCategoryParams(model),
	)
}

func (service *CategoryService) CategoryDelete(
	slug string,
) error {
	return service.repository.CategoryDelete(slug)
}
