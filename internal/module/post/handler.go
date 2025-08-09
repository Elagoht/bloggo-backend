package post

import (
	"bloggo/internal/module/post/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type PostHandler struct {
	service PostService
}

func NewPostHandler(service PostService) PostHandler {
	return PostHandler{
		service,
	}
}

func (handler *PostHandler) ListPosts(
	writer http.ResponseWriter,
	request *http.Request,
) {
	details, err := handler.service.GetPostList()
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(details)
}

func (handler *PostHandler) GetPostBySlug(
	writer http.ResponseWriter,
	request *http.Request,
) {
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	details, err := handler.service.GetPostBySlug(slug)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(details)
}

func (handler *PostHandler) CreatePostWithFirstVersion(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	body, files, ok := handlers.BindAndValidateMultipart[*models.RequestPostUpsert](
		writer, request, 20<<20,
	)
	if !ok {
		return
	}

	handler.service.CreatePostWithFirstVersion(body, files["cover"], userId)

	json.NewEncoder(writer).Encode(body)
}
