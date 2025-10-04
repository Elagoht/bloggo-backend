package config

import (
	"bloggo/internal/utils/validate"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 int    `validate:"required,port"`
	JWTSecret            string `validate:"required,min=32"`
	AccessTokenDuration  int    `validate:"required"`
	RefreshTokenDuration int    `validate:"required"`
	GeminiAPIKey         string
	TrustedFrontendKey   string `validate:"required,min=32"`
}

var (
	instance Config
	once     sync.Once
	loadErr  error
)

// MustLoad loads configuration from environment variables and returns an error if it fails
func MustLoad() error {
	once.Do(func() {
		instance, loadErr = load()
	})
	return loadErr
}

// Get returns the singleton Config instance
func Get() Config {
	return instance
}

func IsGeminiEnabled() bool {
	return Get().GeminiAPIKey != ""
}

func load() (Config, error) {
	// Load .env file if it exists (optional - for local development)
	_ = godotenv.Load()

	// Get port from environment variable, default to 8723
	port, err := getEnvAsInt("PORT", 8723)
	if err != nil {
		return Config{}, err
	}

	// Get JWT secret - REQUIRED
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	// Get access token duration - default to 900 seconds (15 minutes)
	accessTokenDuration, err := getEnvAsInt("ACCESS_TOKEN_DURATION", 900)
	if err != nil {
		return Config{}, err
	}

	// Get refresh token duration - default to 604800 seconds (7 days)
	refreshTokenDuration, err := getEnvAsInt("REFRESH_TOKEN_DURATION", 604800)
	if err != nil {
		return Config{}, err
	}

	// Get Gemini API key - optional
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	// Get trusted frontend key - REQUIRED
	trustedFrontendKey := os.Getenv("TRUSTED_FRONTEND_KEY")
	if trustedFrontendKey == "" {
		return Config{}, fmt.Errorf("TRUSTED_FRONTEND_KEY environment variable is required")
	}

	result := Config{
		Port:                 port,
		JWTSecret:            jwtSecret,
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
		GeminiAPIKey:         geminiAPIKey,
		TrustedFrontendKey:   trustedFrontendKey,
	}

	// Validate configuration
	if err := validate.GetValidator().Struct(result); err != nil {
		return Config{}, fmt.Errorf("configuration is not valid: %w", err)
	}

	return result, nil
}

// getEnvAsInt reads an environment variable as an integer with a default fallback
func getEnvAsInt(key string, defaultValue int) (int, error) {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %s", key, valueStr)
	}

	return value, nil
}
