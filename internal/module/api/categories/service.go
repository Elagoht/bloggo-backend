package categories

import (
	"bloggo/internal/module/api/categories/models"
)

type CategoriesAPIService struct {
	repository CategoriesAPIRepository
}

func NewCategoriesAPIService(repository CategoriesAPIRepository) CategoriesAPIService {
	return CategoriesAPIService{repository}
}

func (service *CategoriesAPIService) GetAllCategories() (*models.APICategoriesResponse, error) {
	return service.repository.GetAllCategories()
}

func (service *CategoriesAPIService) GetCategoryBySlug(slug string) (*models.APICategoryDetails, error) {
	return service.repository.GetCategoryBySlug(slug)
}
