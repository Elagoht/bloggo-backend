package storage

import (
	"bloggo/internal/utils/handlers"
	"net/http"
	"path/filepath"
)

type StorageHandler struct{}

func NewStorageHandler() StorageHandler {
	return StorageHandler{}
}

func (handler *StorageHandler) ServeUserAvatars(
	writer http.ResponseWriter,
	request *http.Request,
) {
	handler.serveImage(writer, request, "uploads", "users", "avatars")
}

func (handler *StorageHandler) ServePostCovers(
	writer http.ResponseWriter,
	request *http.Request,
) {
	handler.serveImage(writer, request, "uploads", "posts", "versions", "covers")
}

func (handler *StorageHandler) serveImage(
	writer http.ResponseWriter,
	request *http.Request,
	pathComponents ...string,
) {
	imageId, ok := handlers.GetParam[string](writer, request, "imageId")
	if !ok {
		return
	}

	// Build path with imageId and .webp extension
	fullPath := append(pathComponents, imageId+".webp")
	imagePath := filepath.Join(fullPath...)

	// Add caching headers
	writer.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day
	// Set response type
	writer.Header().Set("Content-Type", "image/webp")

	// Serve the file
	http.ServeFile(writer, request, imagePath)
}
