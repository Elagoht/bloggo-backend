package category

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/category/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"bloggo/internal/utils/schemas/responses"
)

type CategoryService struct {
	repository  CategoryRepository
	permissions permissions.Store
}

func NewCategoryService(repository CategoryRepository, permissions permissions.Store) CategoryService {
	return CategoryService{
		repository,
		permissions,
	}
}

func (service *CategoryService) CategoryCreate(
	model *models.RequestCategoryCreate,
	userRoleId int64,
) (*responses.ResponseCreated, error) {
	// Check if user has permission to manage categories (editors/admins)
	hasPermission := service.permissions.HasPermission(userRoleId, "category:manage")
	if !hasPermission {
		return nil, apierrors.ErrForbidden
	}

	id, err := service.repository.CategoryCreate(
		models.ToCreateCategoryParams(model),
	)
	if err != nil {
		return nil, err
	}

	return &responses.ResponseCreated{
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
	userRoleId int64,
) error {
	// Check if user has permission to manage categories (editors/admins)
	hasPermission := service.permissions.HasPermission(userRoleId, "category:manage")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	return service.repository.CategoryUpdate(
		slug,
		models.ToUpdateCategoryParams(model),
	)
}

func (service *CategoryService) CategoryDelete(
	slug string,
	userRoleId int64,
) error {
	// Check if user has permission to manage categories (editors/admins)
	hasPermission := service.permissions.HasPermission(userRoleId, "category:manage")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	return service.repository.CategoryDelete(slug)
}
