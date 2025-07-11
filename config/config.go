package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Services map[string]ServiceTypeConfig `yaml:"services"`
}

type ServiceTypeConfig struct {
	Environments map[string]EnvironmentConfig `yaml:"environments"`
}

type EnvironmentConfig struct {
	ProjectID string `yaml:"project_id"`
	Region    string `yaml:"region,omitempty"`
	Cluster   string `yaml:"cluster,omitempty"`
}

// LoadConfig reads and parses the YAML configuration file.
// It takes a filepath argument but currently ignores it.
func LoadConfig(filepathArgument string) (*Config, error) {
	var configFilePath string

	if filepathArgument != "" {
		// Use the provided filepath argument if it's not empty
		configFilePath = filepathArgument
	} else {
		// Get the path of the currently running executable
		execPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("error finding executable path: %w", err)
		}

		// Get the directory containing the executable
		execDir := filepath.Dir(execPath)

		// Construct the path to the configuration file in the executable's directory.
		configFilePath = filepath.Join(execDir, ".gcp-launch.yaml")
	}

	// Read the entire content of the YAML file
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		// Return an error if the file cannot be read (e.g., not found, permissions)
		return nil, fmt.Errorf("error reading config file '%s': %w", configFilePath, err)
	}

	// Create an empty Config struct instance to populate
	var cfg Config

	// Parse the YAML content read from the file into the cfg struct
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		// Return an error if the YAML content is invalid or doesn't match the struct
		return nil, fmt.Errorf("error parsing config file '%s': %w", configFilePath, err)
	}

	// If everything is successful, return a pointer to the populated struct and a nil error
	return &cfg, nil
}
