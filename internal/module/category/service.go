package category

import "bloggo/internal/module/category/models"

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
