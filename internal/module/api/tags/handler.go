package tags

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type TagsAPIHandler struct {
	service TagsAPIService
}

func NewTagsAPIHandler(service TagsAPIService) TagsAPIHandler {
	return TagsAPIHandler{service}
}

func (h *TagsAPIHandler) ListTags(writer http.ResponseWriter, request *http.Request) {
	response, err := h.service.GetAllTags()
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(response)
}

func (h *TagsAPIHandler) GetTagBySlug(writer http.ResponseWriter, request *http.Request) {
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	tag, err := h.service.GetTagBySlug(slug)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(tag)
}
