package config

import (
	"encoding/json"
	"os"
)

// Config holds user configuration
type Config struct {
	XDMPath     string `json:"xdm_path"`
	DownloadDir string `json:"download_dir"`
}

const configFileName = "config.json"

// LoadConfig reads config from file or returns default
func LoadConfig() (*Config, error) {
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return &Config{
			DownloadDir: "Downloads",
		}, nil
	}

	file, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig writes config to file
func SaveConfig(cfg *Config) error {
	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}
