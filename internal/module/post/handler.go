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

	body, ok := handlers.BindAndValidateMultipart[*models.RequestPostUpsert](
		writer,
		request,
		10<<20,
	)
	if !ok {
		return
	}

	createdId, err := handler.service.CreatePostWithFirstVersion(
		body,
		userId,
	)
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

func (handler *PostHandler) CreateVersionFromLatest(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	postId, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	createdId, err := handler.service.CreateVersionFromLatest(postId, userId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(createdId)
}

func (handler *PostHandler) UpdateUnsubmittedOwnVersion(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	postId, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}
	versionId, ok := handlers.GetParam[int64](writer, request, "versionId")
	if !ok {
		return
	}

	body, ok :=
		handlers.BindAndValidateMultipart[*models.RequestPostUpsert](
			writer,
			request,
			10<<20,
		)
	if !ok {
		return
	}

	if err := handler.service.UpdateUnsubmittedOwnVersion(
		postId,
		versionId,
		userId,
		body,
	); err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			// Custom error messages for this endpoint
			apierrors.ErrPreconditionFailed: {
				Message: "This version has already been submitted and cannot be edited.",
				Status:  http.StatusPreconditionFailed,
			},
			apierrors.ErrForbidden: {
				Message: "You can only edit your own versions. Please create a new version.",
				Status:  http.StatusPreconditionFailed,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *PostHandler) SubmitVersionForReview(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	postId, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}
	versionId, ok := handlers.GetParam[int64](writer, request, "versionId")
	if !ok {
		return
	}

	if err := handler.service.SubmitVersionForReview(
		postId,
		versionId,
		userId,
	); err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrPreconditionFailed: {
				Message: "Only draft versions can be submitted for review.",
				Status:  http.StatusPreconditionFailed,
			},
			apierrors.ErrForbidden: {
				Message: "You can only submit your own versions for review.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *PostHandler) ApproveVersion(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	postId, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}
	versionId, ok := handlers.GetParam[int64](writer, request, "versionId")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[*models.RequestPostStatusModerate](
		writer,
		request,
	)
	if !ok {
		return
	}

	var note *string
	if body.Note != "" {
		note = &body.Note
	}

	if err := handler.service.ApproveVersion(
		postId,
		versionId,
		userId,
		note,
	); err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrPreconditionFailed: {
				Message: "Drafts and published versions cannot be approved.",
				Status:  http.StatusPreconditionFailed,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *PostHandler) RejectVersion(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	postId, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}
	versionId, ok := handlers.GetParam[int64](writer, request, "versionId")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[*models.RequestPostStatusModerate](
		writer,
		request,
	)
	if !ok {
		return
	}

	var note *string
	if body.Note != "" {
		note = &body.Note
	}

	if err := handler.service.RejectVersion(
		postId,
		versionId,
		userId,
		note,
	); err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrPreconditionFailed: {
				Message: "Drafts and published versions cannot be rejected.",
				Status:  http.StatusPreconditionFailed,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
