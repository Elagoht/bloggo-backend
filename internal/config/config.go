package config

import (
	"bloggo/internal/utils/cryptography"
	"bloggo/internal/utils/validate"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Port                 int    `json:"port" validate:"required,port"`
	JWTSecret            string `json:"JWTSecret" validate:"required,min=32,max=32"`
	AccessTokenDuration  int    `json:"accessTokenDuration" validate:"required"`
	RefreshTokenDuration int    `json:"refreshTokenDuration" validate:"required"`
	GeminiAPIKey         string `json:"geminiApiKey"`
	TrustedFrontendKey   string `json:"trustedFrontendKey" validate:"required,min=32,max=32"`
}

var (
	instance   Config
	once       sync.Once
	configFile = "bloggo-config.json"
)

func (conf Config) Save(file string) {
	// Ensure the directory exists
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write to file
	if err := os.WriteFile(file, data, 0644); err != nil {
		log.Fatal(err)
	}
}

// Get returns the singleton Config instance, loading it from file if necessary.
func Get() Config {
	once.Do(func() {
		instance = load(configFile)
	})
	return instance
}

func IsGeminiEnabled() bool {
	return Get().GeminiAPIKey != ""
}

func load(file string) Config {
	// If a config file doesn't exist, generate one.
	if _, err := os.Stat(file); os.IsNotExist(err) {
		generateConfig().Save(file)
	}

	// Read config
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Bind file content to config
	result := Config{}
	if err := json.Unmarshal(content, &result); err != nil {
		log.Fatal("Error while parsing configuration.")
	}

	// Validate binded data
	err = validate.GetValidator().Struct(result)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Configuration loaded but is not valid.")
	}
	return result
}

func generateConfig() *Config {
	secret, err := cryptography.GenerateRandomHS256Secret()
	if err != nil {
		log.Fatal("Couldn't generate secret key.")
	}

	trustedFrontendKey, err := cryptography.GenerateRandomHS256Secret()
	if err != nil {
		log.Fatal("Couldn't generate trusted frontend key.")
	}

	return &Config{
		Port:                 8723,               // Default port
		JWTSecret:            secret,             // Random secret key per distributed instance
		AccessTokenDuration:  60 * 15,            // 15 minutes for access token
		RefreshTokenDuration: 60 * 60 * 24 * 7,   // Defaults 7 days for refresh token
		GeminiAPIKey:         "",                 // Empty by default - users can add their key
		TrustedFrontendKey:   trustedFrontendKey, // Random key for trusted frontend requests
	}
}
