package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Port int `json:"port" validate:"required"`
}

func (conf Config) Save(file string) {
	// Ensure the directory exists
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		panic(err)
	}

	// Write to file
	if err := os.WriteFile(file, data, 0644); err != nil {
		panic(err)
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
		panic(err)
	}

	// Convert to object and return
	result := Config{}
	if err := json.Unmarshal(content, &result); err != nil {
		panic("Error while parsing configuration.")
	}
	return result
}
