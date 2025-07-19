package validate

import (
	"path/filepath"
	"regexp"
	"strconv"
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
				panic(err)
			}
		}
	})
	return validatorInstance
}

var customValidator = map[string]func(validator.FieldLevel) bool{
	"port":     PortStringValidator,
	"safePath": SafePathValidator,
}

// Checks if a string is a valid port number (1025-65535)
func PortStringValidator(fl validator.FieldLevel) bool {
	portStr := fl.Field().String()
	if portStr == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^\d+$`, portStr)
	if !matched {
		return false
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}
	return port >= 1025 && port <= 65535
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
