package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if config.CopyFiles == nil {
		t.Error("DefaultConfig() CopyFiles should not be nil")
	}

	// Default config should now have empty copy_files array
	if len(config.CopyFiles) != 0 {
		t.Errorf("DefaultConfig() should have empty copy_files, got %d files", len(config.CopyFiles))
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				CopyFiles: []string{".env", "docker-compose.yml"},
			},
			wantErr: false,
		},
		{
			name: "nil copy_files",
			config: &Config{
				CopyFiles: nil,
			},
			wantErr: true,
		},
		{
			name: "empty copy_files",
			config: &Config{
				CopyFiles: []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindConfigFile(t *testing.T) {
	tempDir := t.TempDir()

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(tempDir, ConfigFileName)
	if err := os.WriteFile(configPath, []byte(`{"copy_files": []}`), 0o644); err != nil {
		t.Fatal(err)
	}

	foundPath, err := FindConfigFile(subDir)
	if err != nil {
		t.Fatalf("FindConfigFile() error = %v", err)
	}

	if foundPath != configPath {
		t.Errorf("FindConfigFile() = %v, want %v", foundPath, configPath)
	}

	emptyDir := t.TempDir()
	_, err = FindConfigFile(emptyDir)
	if err == nil {
		t.Error("FindConfigFile() should return error for non-existent config")
	}
}

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		content  string
		wantErr  bool
		validate func(*Config) bool
	}{
		{
			name:    "valid config",
			content: `{"copy_files": [".env", "docker-compose.yml"]}`,
			wantErr: false,
			validate: func(c *Config) bool {
				return len(c.CopyFiles) == 2 && c.CopyFiles[0] == ".env"
			},
		},
		{
			name:    "invalid json",
			content: `{"copy_files": [".env"`,
			wantErr: true,
		},
		{
			name:    "invalid config structure",
			content: `{"copy_files": null}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tempDir, "test_"+tt.name+".json")
			if err := os.WriteFile(configPath, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			config, err := LoadConfig(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil && !tt.validate(config) {
				t.Error("LoadConfig() validation failed")
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.json")

	config := &Config{
		CopyFiles: []string{".env", "docker-compose.yml"},
	}

	err := SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if len(loadedConfig.CopyFiles) != len(config.CopyFiles) {
		t.Error("Saved config does not match original")
	}
}

func TestCreateDefaultConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "new_config.json")

	err := CreateDefaultConfigFile(configPath)
	if err != nil {
		t.Fatalf("CreateDefaultConfigFile() error = %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Default config file was not created")
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load created default config: %v", err)
	}

	if config.CopyFiles == nil {
		t.Error("Default config CopyFiles should not be nil")
	}

	if len(config.CopyFiles) != 0 {
		t.Errorf("Default config should have empty copy_files, got %d files", len(config.CopyFiles))
	}

	err = CreateDefaultConfigFile(configPath)
	if err == nil {
		t.Error("CreateDefaultConfigFile() should return error when file already exists")
	}
}
