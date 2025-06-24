package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const ConfigFileName = "sproutee.json"

type Config struct {
	CopyFiles []string `json:"copy_files"`
}

func DefaultConfig() *Config {
	return &Config{
		CopyFiles: []string{},
	}
}

func (c *Config) Validate() error {
	if c.CopyFiles == nil {
		return fmt.Errorf("copy_files field is required")
	}
	return nil
}

func FindConfigFile(startDir string) (string, error) {
	currentDir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	for {
		configPath := filepath.Join(currentDir, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("configuration file '%s' not found", ConfigFileName)
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func LoadConfigFromCurrentDir() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	configPath, err := FindConfigFile(wd)
	if err != nil {
		return nil, err
	}

	return LoadConfig(configPath)
}

func SaveConfig(config *Config, configPath string) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func CreateDefaultConfigFile(configPath string) error {
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("configuration file already exists: %s", configPath)
	}

	defaultConfig := DefaultConfig()
	return SaveConfig(defaultConfig, configPath)
}
