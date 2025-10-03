package category

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/ai"
	"bloggo/internal/module/audit"
	auditmodels "bloggo/internal/module/audit/models"
	"bloggo/internal/module/category/models"
	"bloggo/internal/module/webhook"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"bloggo/internal/utils/schemas/responses"
)

type CategoryService struct {
	repository  CategoryRepository
	permissions permissions.Store
	aiService   ai.AIService
}

func NewCategoryService(repository CategoryRepository, permissions permissions.Store, aiService ai.AIService) CategoryService {
	return CategoryService{
		repository,
		permissions,
		aiService,
	}
}

func (service *CategoryService) CategoryCreate(
	model *models.RequestCategoryCreate,
	userRoleId int64,
	userId int64,
) (*responses.ResponseCreated, error) {
	// Check if user has permission to create categories
	hasPermission := service.permissions.HasPermission(userRoleId, "category:create")
	if !hasPermission {
		return nil, apierrors.ErrForbidden
	}

	params := models.ToCreateCategoryParams(model)
	id, err := service.repository.CategoryCreate(params)
	if err != nil {
		return nil, err
	}

	// Log the action
	audit.LogCategoryAction(&userId, id, auditmodels.ActionCategoryCreated)

	// Trigger webhook
	go func() {
		webhook.TriggerCategoryCreated(id, params.Slug, map[string]interface{}{
			"name":        params.Name,
			"slug":        params.Slug,
			"spot":        params.Spot,
			"description": params.Description,
		})
	}()

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
) (*responses.PaginatedResponse[models.ResponseCategoryCard], error) {
	// Get the categories data
	categories, err := service.repository.GetCategories(pagination, search)
	if err != nil {
		return nil, err
	}

	// Get the total count with same filters
	total, err := service.repository.GetCategoriesCount(search)
	if err != nil {
		return nil, err
	}

	// Set default values for page and take if they're nil
	page := 1
	if pagination.Page != nil {
		page = *pagination.Page
	}

	take := 12 // default take value
	if pagination.Take != nil {
		take = *pagination.Take
	}

	return &responses.PaginatedResponse[models.ResponseCategoryCard]{
		Data:  categories,
		Page:  page,
		Take:  take,
		Total: total,
	}, nil
}

func (service *CategoryService) GetCategoryList() ([]models.ResponseCategoryListItem, error) {
	return service.repository.GetCategoryList()
}

func (service *CategoryService) CategoryUpdate(
	slug string,
	model *models.RequestCategoryUpdate,
	userRoleId int64,
	userId int64,
) error {
	// Check if user has permission to update categories
	hasPermission := service.permissions.HasPermission(userRoleId, "category:update")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	// First get the category ID to log it
	category, err := service.repository.GetCategoryBySlug(slug)
	if err != nil {
		return err
	}

	params := models.ToUpdateCategoryParams(model)
	err = service.repository.CategoryUpdate(slug, params)
	if err != nil {
		return err
	}

	// Log the action
	audit.LogCategoryAction(&userId, category.Id, auditmodels.ActionCategoryUpdated)

	// Trigger webhook with updated slug
	newSlug := slug
	if params.Slug != nil {
		newSlug = *params.Slug
	}
	go func() {
		webhook.TriggerCategoryUpdated(category.Id, newSlug, map[string]interface{}{
			"name":        model.Name,
			"slug":        newSlug,
			"spot":        model.Spot,
			"description": model.Description,
		})
	}()

	return nil
}

func (service *CategoryService) CategoryDelete(
	slug string,
	userRoleId int64,
	userId int64,
) error {
	// Check if user has permission to delete categories
	hasPermission := service.permissions.HasPermission(userRoleId, "category:delete")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	// First get the category details for audit logging
	category, err := service.repository.GetCategoryBySlug(slug)
	if err != nil {
		return err
	}

	// Soft delete the category (posts will keep reference but show as archived)
	err = service.repository.CategoryDelete(slug)
	if err != nil {
		return err
	}

	// Log the action
	audit.LogCategoryAction(&userId, category.Id, auditmodels.ActionCategoryDeleted)

	// Trigger webhook
	go func() {
		webhook.TriggerCategoryDeleted(category.Id, category.Slug)
	}()

	return nil
}

func (service *CategoryService) GenerativeFill(
	categoryName string,
	userRoleId int64,
) (*models.ResponseCategoryGenerativeFill, error) {
	// Check if user has permission to create categories
	hasPermission := service.permissions.HasPermission(userRoleId, "category:create")
	if !hasPermission {
		return nil, apierrors.ErrForbidden
	}

	result, err := service.aiService.GenerateCategoryMetadata(categoryName)
	if err != nil {
		return nil, err
	}

	return &models.ResponseCategoryGenerativeFill{
		Spot:        result.Spot,
		Description: result.Description,
	}, nil
}
