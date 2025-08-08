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
	imageId, ok := handlers.GetParam[string](writer, request, "imageId")
	if !ok {
		return
	}

	avatarPath := filepath.Join(
		"uploads",
		"users",
		"avatars",
		imageId+".webp",
	)

	// Add caching headers
	writer.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day
	// Set response type
	writer.Header().Set("Content-Type", "image/webp")

	// Serve the file
	http.ServeFile(writer, request, avatarPath)
}
