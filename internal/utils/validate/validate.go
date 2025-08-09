package validate

import (
	"log"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validatorInstance *validator.Validate
	once              sync.Once
)

// Get singleton instance
func GetValidator() *validator.Validate {
	once.Do(func() {
		validatorInstance = validator.New()
		for key, validateFunction := range customValidator {
			err := validatorInstance.RegisterValidation(key, validateFunction)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	return validatorInstance
}

var customValidator = map[string]func(validator.FieldLevel) bool{
	"port":     PortValidator,
	"safePath": SafePathValidator,
	"file":     FileValidator,
}

// Checks if a number is a valid port number (80, 443 or in range 1025-65535)
func PortValidator(fieldLevel validator.FieldLevel) bool {
	port := fieldLevel.Field().Int()
	return port == 80 || port == 443 || (port >= 1024 && port <= 65535)
}

// Checks if a path is valid and does not contain path traversal
func SafePathValidator(fieldLevel validator.FieldLevel) bool {
	path := fieldLevel.Field().String()
	if path == "" {
		return false
	}
	cleaned := filepath.Clean(path)
	return !strings.Contains(cleaned, "..")
}

// Check files
func FileValidator(fieldLevel validator.FieldLevel) bool {
	// Field type can be string (file name) or multipart.FileHeader
	if fieldLevel.Field().Kind() == reflect.String {
		return fieldLevel.Field().String() != ""
	}

	if fileHeader, ok := fieldLevel.Field().Interface().(*multipart.FileHeader); ok {
		return fileHeader != nil && fileHeader.Size > 0
	}

	return false
}
