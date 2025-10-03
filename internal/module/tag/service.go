package tag

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/audit"
	auditmodels "bloggo/internal/module/audit/models"
	"bloggo/internal/module/tag/models"
	"bloggo/internal/module/webhook"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"bloggo/internal/utils/schemas/responses"
)

type TagService struct {
	repository  TagRepository
	permissions permissions.Store
}

func NewTagService(repository TagRepository, permissions permissions.Store) TagService {
	return TagService{
		repository,
		permissions,
	}
}

func (service *TagService) TagCreate(
	model *models.RequestTagCreate,
	userRoleId int64,
	userId int64,
) (*responses.ResponseCreated, error) {
	// Check if user has permission to create tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:create")
	if !hasPermission {
		return nil, apierrors.ErrForbidden
	}

	params := models.ToCreateTagParams(model)
	id, err := service.repository.TagCreate(params)
	if err != nil {
		return nil, err
	}

	// Log the action
	audit.LogTagAction(&userId, id, auditmodels.ActionTagCreated)

	// Trigger webhook
	go func() { webhook.TriggerTagCreated(id, params.Slug, map[string]interface{}{"name": params.Name, "slug": params.Slug}) }()

	return &responses.ResponseCreated{
		Id: id,
	}, nil
}

func (service *TagService) GetTagBySlug(
	slug string,
) (*models.ResponseTagDetails, error) {
	return service.repository.GetTagBySlug(slug)
}

func (service *TagService) GetTags(
	pagination *pagination.PaginationOptions,
	search *filter.SearchOptions,
) (*responses.PaginatedResponse[models.ResponseTagCard], error) {
	tags, total, err := service.repository.GetTags(pagination, search)
	if err != nil {
		return nil, err
	}

	page := 1
	take := 24
	if pagination.Page != nil {
		page = *pagination.Page
	}
	if pagination.Take != nil {
		take = *pagination.Take
	}

	return &responses.PaginatedResponse[models.ResponseTagCard]{
		Data:  tags,
		Page:  page,
		Take:  take,
		Total: total,
	}, nil
}

func (service *TagService) GetTagList() ([]models.ResponseTagListItem, error) {
	return service.repository.GetTagList()
}

func (service *TagService) TagUpdate(
	slug string,
	model *models.RequestTagUpdate,
	userRoleId int64,
	userId int64,
) error {
	// Check if user has permission to update tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:update")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	// First get the tag ID to log it
	tag, err := service.repository.GetTagBySlug(slug)
	if err != nil {
		return err
	}

	params := models.ToUpdateTagParams(model)
	err = service.repository.TagUpdate(slug, params)
	if err != nil {
		return err
	}

	// Log the action
	audit.LogTagAction(&userId, tag.Id, auditmodels.ActionTagUpdated)

	// Trigger webhook with updated slug
	newSlug := slug
	if params.Slug != nil {
		newSlug = *params.Slug
	}
	go func() {
		webhook.TriggerTagUpdated(tag.Id, newSlug, map[string]interface{}{"name": model.Name, "slug": newSlug})
	}()

	return nil
}

func (service *TagService) TagDelete(
	slug string,
	userRoleId int64,
	userId int64,
) error {
	// Check if user has permission to delete tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:delete")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	// First get the tag ID to log it
	tag, err := service.repository.GetTagBySlug(slug)
	if err != nil {
		return err
	}

	err = service.repository.TagDelete(slug)
	if err != nil {
		return err
	}

	// Log the action
	audit.LogTagAction(&userId, tag.Id, auditmodels.ActionTagDeleted)

	// Trigger webhook
	go func() { webhook.TriggerTagDeleted(tag.Id, tag.Slug) }()

	return nil
}
