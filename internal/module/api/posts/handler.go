package posts

import (
	"bloggo/internal/module/api/posts/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"encoding/json"
	"net/http"
)

type PostsAPIHandler struct {
	service PostsAPIService
}

func NewPostsAPIHandler(service PostsAPIService) PostsAPIHandler {
	return PostsAPIHandler{service}
}

func (h *PostsAPIHandler) ListPublishedPosts(writer http.ResponseWriter, request *http.Request) {
	// Get pagination parameters
	page, _ := handlers.GetQuery[int](writer, request, "page")
	limit, _ := handlers.GetQuery[int](writer, request, "limit")

	// Get filter parameters
	categorySlug := request.URL.Query().Get("category")
	tagSlug := request.URL.Query().Get("tag")
	authorId := request.URL.Query().Get("author")
	search := request.URL.Query().Get("search")

	var categoryPtr, tagPtr, authorPtr, searchPtr *string
	if categorySlug != "" {
		categoryPtr = &categorySlug
	}
	if tagSlug != "" {
		tagPtr = &tagSlug
	}
	if authorId != "" {
		authorPtr = &authorId
	}
	if search != "" {
		searchPtr = &search
	}

	response, err := h.service.GetPublishedPosts(page, limit, categoryPtr, tagPtr, authorPtr, searchPtr)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(response)
}

func (h *PostsAPIHandler) GetPublishedPostBySlug(writer http.ResponseWriter, request *http.Request) {
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	post, err := h.service.GetPublishedPostBySlug(slug)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(post)
}

func (h *PostsAPIHandler) TrackPostView(writer http.ResponseWriter, request *http.Request) {
	slug, ok := handlers.GetParam[string](writer, request, "slug")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[*models.APITrackViewRequest](writer, request)
	if !ok {
		return
	}

	err := h.service.TrackView(slug, body.UserAgent)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (h *PostsAPIHandler) GetAllViewCounts(writer http.ResponseWriter, request *http.Request) {
	viewCounts, err := h.service.GetAllViewCounts()
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(viewCounts)
}
