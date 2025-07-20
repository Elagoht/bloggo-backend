package validate

import (
	"log"
	"path/filepath"
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
