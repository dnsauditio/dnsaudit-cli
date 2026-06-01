package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey string `json:"api_key"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "dnsaudit")
	
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// Load reads the API key from env or config file
func Load() (string, error) {
	// 1. Check environment variable first
	apiKey := os.Getenv("DNSAUDIT_API_KEY")
	if apiKey != "" {
		return apiKey, nil
	}

	// 2. Fallback to config file
	configPath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No config file yet
		}
		return "", err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	return cfg.APIKey, nil
}

// Save writes the API key to the config file
func Save(apiKey string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	cfg := Config{APIKey: apiKey}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return err
	}

	fmt.Printf("API Key saved to %s\n", configPath)
	return nil
}
