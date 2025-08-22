package category

import (
	"bloggo/internal/module/category/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/handlers"
	"bloggo/internal/utils/pagination"
	"encoding/json"
	"net/http"
)

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) CategoryHandler {
	return CategoryHandler{
		service,
	}
}

func (handler *CategoryHandler) CategoryCreate(
	writer http.ResponseWriter,
	request *http.Request,
) {
	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[models.RequestCategoryCreate](writer, request)
	if !ok {
		return
	}

	response, err := handler.service.CategoryCreate(&body, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage categories.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(response)
}

func (handler *CategoryHandler) GetCategoryBySlug(
	writer http.ResponseWriter,
	request *http.Request,
) {
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	details, err := handler.service.GetCategoryBySlug(slug)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(details)
}

func (handler *CategoryHandler) GetCategories(
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

func (handler *CategoryHandler) GetCategoryList(
	writer http.ResponseWriter,
	request *http.Request,
) {
	categories, err := handler.service.GetCategoryList()
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(categories)
}

func (handler *CategoryHandler) CategoryUpdate(
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

	body, ok := handlers.BindAndValidate[models.RequestCategoryUpdate](writer, request)
	if !ok {
		return
	}

	err := handler.service.CategoryUpdate(slug, &body, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage categories.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *CategoryHandler) CategoryDelete(
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

	err := handler.service.CategoryDelete(slug, roleId)
	if err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrForbidden: {
				Message: "Only editors and admins can manage categories.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
