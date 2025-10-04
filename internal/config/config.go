package config

import (
	"bloggo/internal/utils/validate"
	"fmt"
	"log"
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
)

// Get returns the singleton Config instance, loading it from environment variables.
func Get() Config {
	once.Do(func() {
		instance = load()
	})
	return instance
}

func IsGeminiEnabled() bool {
	return Get().GeminiAPIKey != ""
}

func load() Config {
	// Load .env file if it exists (optional - for local development)
	_ = godotenv.Load()

	// Get port from environment variable, default to 8723
	port := getEnvAsInt("PORT", 8723)

	// Get JWT secret - REQUIRED
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Get access token duration - default to 900 seconds (15 minutes)
	accessTokenDuration := getEnvAsInt("ACCESS_TOKEN_DURATION", 900)

	// Get refresh token duration - default to 604800 seconds (7 days)
	refreshTokenDuration := getEnvAsInt("REFRESH_TOKEN_DURATION", 604800)

	// Get Gemini API key - optional
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	// Get trusted frontend key - REQUIRED
	trustedFrontendKey := os.Getenv("TRUSTED_FRONTEND_KEY")
	if trustedFrontendKey == "" {
		log.Fatal("TRUSTED_FRONTEND_KEY environment variable is required")
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
	err := validate.GetValidator().Struct(result)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Configuration is not valid")
	}

	return result
}

// getEnvAsInt reads an environment variable as an integer with a default fallback
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", key, valueStr)
	}

	return value
}
