package validatefile

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type ImageValidator struct {
	MaxSize int64
}

func NewImageValidator(maxSize int64) FileValidator {
	return &ImageValidator{
		MaxSize: maxSize,
	}
}

func (validator *ImageValidator) Validate(
	file multipart.File,
	header *multipart.FileHeader,
) error {
	// Check size
	if header.Size > validator.MaxSize {
		return fmt.Errorf(
			"file is too large: max %d bytes allowed",
			validator.MaxSize,
		)
	}

	// Read first 512 byte to get mimetype
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Check file type from read buffer
	mimeType := http.DetectContentType(buffer)
	if !strings.HasPrefix(mimeType, "image/") {
		return fmt.Errorf("file is not an image")
	}

	// Reset the cursor
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to rewind file: %w", err)
	}

	return nil
}
