package post

import (
	"bloggo/internal/module/post/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/handlers"
	"bloggo/internal/utils/pagination"
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
	// Get pagination options
	paginationOptions, ok := pagination.GetPaginationOptions(
		writer,
		request,
		[]string{"title", "created_at", "updated_at", "read_count"},
	)
	if !ok {
		return
	}

	// Get search options
	searchOptions, ok := filter.GetSearchOptions(writer, request)
	if !ok {
		return
	}

	// Get additional filters
	status, okStatus := handlers.GetQuery[int](writer, request, "status")
	categoryId, okCategoryId := handlers.GetQuery[int64](writer, request, "categoryId")
	authorId, okAuthorId := handlers.GetQuery[int64](writer, request, "authorId")

	var statusPtr *int
	var categoryIdPtr *int64
	var authorIdPtr *int64

	if okStatus {
		statusPtr = &status
	}
	if okCategoryId {
		categoryIdPtr = &categoryId
	}
	if okAuthorId {
		authorIdPtr = &authorId
	}

	filters := &models.RequestPostFilters{
		Page:       paginationOptions.Page,
		Take:       paginationOptions.Take,
		Order:      paginationOptions.OrderBy,
		Dir:        paginationOptions.Direction,
		Q:          searchOptions.Q,
		Status:     statusPtr,
		CategoryId: categoryIdPtr,
		AuthorId:   authorIdPtr,
	}

	details, err := handler.service.GetPostListPaginated(filters)
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

func (handler *PostHandler) DeleteVersionById(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	roleId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenRoleId)
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

	if err := handler.service.DeleteVersionById(
		postId,
		versionId,
		userId,
		roleId,
	); err != nil {
		apierrors.MapErrors(err, writer, apierrors.HTTPErrorMapping{
			apierrors.ErrPreconditionFailed: {
				Message: "Only draft, pending, or rejected versions can be deleted by authors.",
				Status:  http.StatusPreconditionFailed,
			},
			apierrors.ErrForbidden: {
				Message: "You can only delete your own versions unless you have editor permissions.",
				Status:  http.StatusForbidden,
			},
		})
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *PostHandler) TrackView(
	writer http.ResponseWriter,
	request *http.Request,
) {
	body, ok := handlers.BindAndValidate[*models.RequestTrackView](
		writer,
		request,
	)
	if !ok {
		return
	}

	err := handler.service.TrackView(body)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
