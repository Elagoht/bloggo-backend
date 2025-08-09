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

func (handler *PostHandler) GetPostById(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	details, err := handler.service.GetPostById(id)
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

	createdId, err := handler.service.CreatePostWithFirstVersion(body, files["cover"], userId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(createdId)
}

func (handler *PostHandler) ListPostVersionsGetByPostId(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	versions, err := handler.service.ListPostVersionsGetByPostId(id)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(versions)
}

func (handler *PostHandler) GetPostVersionById(
	writer http.ResponseWriter,
	request *http.Request,
) {
	postId, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}
	versionId, ok := handlers.GetParam[int64](writer, request, "versionId")
	if !ok {
		return
	}

	version, err := handler.service.GetPostVersionById(postId, versionId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(version)
}

func (handler *PostHandler) DeletePostById(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	if err := handler.service.DeletePostById(id); err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
