package config

import (
	"bloggo/internal/utils/validate"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Port int `json:"port" validate:"required,port"`
}

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

func Load(file string) Config {
	// Check if file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// If not exists, create a default one
		defaultConf := Config{Port: 8723}
		defaultConf.Save(file)
	}

	// Read content
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Convert to object and return
	result := Config{}
	if err := json.Unmarshal(content, &result); err != nil {
		log.Fatal("Error while parsing configuration.")
	}

	// Validate unmarshalled data
	err = validate.GetValidator().Struct(result)
	if err != nil {
		log.Fatal("Configuration loaded but is not valid.")
	}

	return result
}
