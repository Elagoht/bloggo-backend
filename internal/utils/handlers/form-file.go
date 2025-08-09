package handlers

import (
	"bloggo/internal/utils/apierrors"
	"mime/multipart"
	"net/http"
)

// GetFormFile retrieves a file from the multipart form data in the request.
func GetFormFile(
	writer http.ResponseWriter,
	request *http.Request,
	fileKey string,
	maxFileSize int64,
) (multipart.File, *multipart.FileHeader, bool) {
	// Parse multipart form data
	err := request.ParseMultipartForm(maxFileSize)
	if err != nil {
		WriteError(
			writer,
			apierrors.NewAPIError("Failed to parse multipart form: "+err.Error(), err),
			http.StatusBadRequest,
		)
		return nil, nil, false
	}

	// Take file from from
	file, fileHeader, err := request.FormFile(fileKey)
	if err != nil {
		if err == http.ErrMissingFile {
			WriteError(
				writer,
				apierrors.NewAPIError("File with key \""+fileKey+"\" not found in form", nil),
				http.StatusBadRequest,
			)
		} else {
			WriteError(
				writer,
				apierrors.NewAPIError("Error retrieving file: "+err.Error(), err),
				http.StatusBadRequest,
			)
		}
		return nil, nil, false
	}

	// Check file size
	if fileHeader.Size > maxFileSize {
		file.Close()
		WriteError(
			writer,
			apierrors.NewAPIError("File size exceeds maximum limit", nil),
			http.StatusRequestEntityTooLarge,
		)
		return nil, nil, false
	}

	return file, fileHeader, true
}
