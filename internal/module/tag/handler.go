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
	body, ok := handlers.BindAndValidate[models.RequestTagCreate](writer, request)
	if !ok {
		return
	}

	response, err := handler.service.TagCreate(&body)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
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
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[models.RequestTagUpdate](writer, request)
	if !ok {
		return
	}

	err := handler.service.TagUpdate(slug, &body)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *TagHandler) TagDelete(
	writer http.ResponseWriter,
	request *http.Request,
) {
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	err := handler.service.TagDelete(slug)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
