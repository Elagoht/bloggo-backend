package tag

import (
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/tag/models"
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
) (*responses.ResponseCreated, error) {
	// Check if user has permission to create tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:create")
	if !hasPermission {
		return nil, apierrors.ErrForbidden
	}

	id, err := service.repository.TagCreate(
		models.ToCreateTagParams(model),
	)
	if err != nil {
		return nil, err
	}

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
	take := 12
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
) error {
	// Check if user has permission to update tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:update")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	return service.repository.TagUpdate(
		slug,
		models.ToUpdateTagParams(model),
	)
}

func (service *TagService) TagDelete(
	slug string,
	userRoleId int64,
) error {
	// Check if user has permission to delete tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:delete")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	return service.repository.TagDelete(slug)
}

// Post-Tag Relationship Methods
func (service *TagService) GetPostTags(postId int64) (*models.ResponsePostTags, error) {
	tags, err := service.repository.GetPostTags(postId)
	if err != nil {
		return nil, err
	}

	return &models.ResponsePostTags{
		PostId: postId,
		Tags:   tags,
	}, nil
}

func (service *TagService) AssignTagsToPost(
	postId int64,
	tagIds []int64,
	userRoleId int64,
) error {
	// Check if user has permission to update tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:update")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	return service.repository.AssignTagsToPost(postId, tagIds)
}

func (service *TagService) RemoveTagFromPost(
	postId int64,
	tagId int64,
	userRoleId int64,
) error {
	// Check if user has permission to update tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:update")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	return service.repository.RemoveTagFromPost(postId, tagId)
}
