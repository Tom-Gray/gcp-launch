package url

import (
	"testing"

	"github.com/tom-gray/gcp-launch/config"
)

func TestGenerateServiceURL(t *testing.T) {
	tests := []struct {
		name        string
		serviceType string
		envConfig   config.EnvironmentConfig
		expectedURL string
		expectError bool
	}{
		{
			name:        "logging service",
			serviceType: "logging",
			envConfig:   config.EnvironmentConfig{ProjectID: "test-project"},
			expectedURL: "https://console.cloud.google.com/logs/viewer?project=test-project",
			expectError: false,
		},
		{
			name:        "cloudrun service with region",
			serviceType: "cloudrun",
			envConfig:   config.EnvironmentConfig{ProjectID: "test-project", Region: "us-central1"},
			expectedURL: "https://console.cloud.google.com/run?project=test-project&region=us-central1",
			expectError: false,
		},
		{
			name:        "gke service with cluster",
			serviceType: "gke",
			envConfig:   config.EnvironmentConfig{ProjectID: "test-project", Cluster: "test-cluster"},
			expectedURL: "https://console.cloud.google.com/kubernetes/workload/overview?inv=1&invt=Ab2VWw&project=test-project",
			expectError: false,
		},
		{
			name:        "spanner service",
			serviceType: "spanner",
			envConfig:   config.EnvironmentConfig{ProjectID: "test-project"},
			expectedURL: "https://console.cloud.google.com/spanner?project=test-project",
			expectError: false,
		},
		{
			name:        "unsupported service type",
			serviceType: "unsupported",
			envConfig:   config.EnvironmentConfig{ProjectID: "test-project"},
			expectedURL: "",
			expectError: true,
		},
		{
			name:        "missing project ID",
			serviceType: "logging",
			envConfig:   config.EnvironmentConfig{ProjectID: ""},
			expectedURL: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := GenerateServiceURL(tt.serviceType, tt.envConfig)
			if (err != nil) != tt.expectError {
				t.Errorf("GenerateServiceURL() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if url != tt.expectedURL {
				t.Errorf("GenerateServiceURL() got URL = %v, want %v", url, tt.expectedURL)
			}
		})
	}
}

func TestGenerateCloudRunURL(t *testing.T) {
	projectID := "my-project"
	region := "asia-east1"
	expected := "https://console.cloud.google.com/run?project=my-project&region=asia-east1"
	actual := GenerateCloudRunURL(projectID, region)
	if actual != expected {
		t.Errorf("GenerateCloudRunURL(%s, %s) = %s; want %s", projectID, region, actual, expected)
	}
}

func TestGenerateGKEURL(t *testing.T) {
	projectID := "my-project"
	cluster := "my-cluster"
	expected := "https://console.cloud.google.com/kubernetes/workload/overview?inv=1&invt=Ab2VWw&project=my-project"
	actual := GenerateGKEURL(projectID, cluster)
	if actual != expected {
		t.Errorf("GenerateGKEURL(%s, %s) = %s; want %s", projectID, cluster, actual, expected)
	}
}
