package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	configContent := `
services:
  logging:
    environments:
      test-env:
        project_id: test-project-id
`
	
	testConfigFile := filepath.Join(tempDir, ".gcp-launch.yaml")
	err := os.WriteFile(testConfigFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Test with explicit file path
	cfg, err := LoadConfig(testConfigFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Services["logging"].Environments["test-env"].ProjectID != "test-project-id" {
		t.Errorf("Expected project_id 'test-project-id', got %s", cfg.Services["logging"].Environments["test-env"].ProjectID)
	}

	// Test with non-existent file
	_, err = LoadConfig("non-existent-file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
