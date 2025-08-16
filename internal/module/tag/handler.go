package tag

import (
	"bloggo/internal/module/tag/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/handlers"
	"bloggo/internal/utils/pagination"
	"encoding/json"
	"net/http"
)

type TagHandler struct {
	service TagService
}

func NewTagHandler(service TagService) TagHandler {
	return TagHandler{
		service,
	}
}

func (handler *TagHandler) TagCreate(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[models.RequestTagCreate](writer, request)
	if !ok {
		return
	}

	response, err := handler.service.TagCreate(&body, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage tags.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(response)
}

func (handler *TagHandler) GetTagBySlug(
	writer http.ResponseWriter,
	request *http.Request,
) {
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	details, err := handler.service.GetTagBySlug(slug)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(details)
}

func (handler *TagHandler) GetCategories(
	writer http.ResponseWriter,
	request *http.Request,
) {
	paginate, ok := pagination.GetPaginationOptions(writer, request, []string{
		"name", "created_at", "updated_at",
	})
	if !ok {
		return
	}

	search, ok := filter.GetSearchOptions(writer, request)
	if !ok {
		return
	}

	categories, err := handler.service.GetCategories(paginate, search)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(categories)
}

func (handler *TagHandler) TagUpdate(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[models.RequestTagUpdate](writer, request)
	if !ok {
		return
	}

	err := handler.service.TagUpdate(slug, &body, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage tags.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *TagHandler) TagDelete(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	err := handler.service.TagDelete(slug, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage tags.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

// Post-Tag Relationship Handlers
func (handler *TagHandler) GetPostTags(
	writer http.ResponseWriter,
	request *http.Request,
) {
	postId, ok := handlers.GetParam[int64](writer, request, "postId")
	if !ok {
		return
	}

	response, err := handler.service.GetPostTags(postId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(response)
}

func (handler *TagHandler) AssignTagsToPost(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	postId, ok := handlers.GetParam[int64](writer, request, "postId")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[models.RequestAssignTagsToPost](writer, request)
	if !ok {
		return
	}

	err := handler.service.AssignTagsToPost(postId, body.TagIds, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage tags.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *TagHandler) RemoveTagFromPost(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	postId, ok := handlers.GetParam[int64](writer, request, "postId")
	if !ok {
		return
	}

	tagId, ok := handlers.GetParam[int64](writer, request, "tagId")
	if !ok {
		return
	}

	err := handler.service.RemoveTagFromPost(postId, tagId, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage tags.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
