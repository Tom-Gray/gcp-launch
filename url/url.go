package url

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/tom-gray/gcp-launch/config"
)

// GenerateServiceURL constructs the appropriate Google Cloud Console URL
// based on the requested service type and environment configuration.
func GenerateServiceURL(serviceType string, envConfig config.EnvironmentConfig) (string, error) {
	if envConfig.ProjectID == "" {
		return "", fmt.Errorf("cannot generate URL: project_id is missing for service type '%s'", serviceType)
	}
	const consoleBaseURL = "https://console.cloud.google.com"
	var url string
	switch serviceType {
	case "logging":
		url = fmt.Sprintf("%s/logs/viewer?project=%s", consoleBaseURL, envConfig.ProjectID)
	case "cloudrun":
		if envConfig.Region != "" {
			url = GenerateCloudRunURL(envConfig.ProjectID, envConfig.Region)
		} else {
			url = fmt.Sprintf("%s/run?project=%s", consoleBaseURL, envConfig.ProjectID)
		}
	case "gke":
		// Use the specific cluster details URL if cluster name is available
		if envConfig.Cluster != "" {
			url = GenerateGKEURL(envConfig.ProjectID, envConfig.Cluster) // Call the new function
		} else {
			// Fallback to the project-level cluster list
			url = fmt.Sprintf("%s/kubernetes/list?project=%s", consoleBaseURL, envConfig.ProjectID)
		}
	case "spanner":
		url = fmt.Sprintf("%s/spanner?project=%s", consoleBaseURL, envConfig.ProjectID)
	default:
		return "", fmt.Errorf("URL generation not supported for service type: '%s'", serviceType)
	}
	return url, nil
}

// GenerateCloudRunURL constructs the Google Cloud Console URL for Cloud Run services
// within a specific project.
func GenerateCloudRunURL(projectID string, region string) string {
	// Cloud Run URLs often include the region for more specific navigation.
	// Format: https://console.cloud.google.com/run?project=<project_id>&region=<region>
	const cloudRunURLFormat = "https://console.cloud.google.com/run?project=%s&region=%s"
	url := fmt.Sprintf(cloudRunURLFormat, projectID, region)
	return url
}

// GenerateGKEURL constructs the Google Cloud Console URL for the GKE workload overview page.
// It uses the format: https://console.cloud.google.com/kubernetes/workload/overview?inv=1&invt=Ab2VWw&project={project_id}
func GenerateGKEURL(projectID string, cluster string) string {
	// Format for workload overview page
	const gkeURLFormat = "https://console.cloud.google.com/kubernetes/workload/overview?inv=1&invt=Ab2VWw&project=%s"
	// Only use project ID for workload overview (cluster parameter is ignored for this URL format)
	url := fmt.Sprintf(gkeURLFormat, projectID)
	return url
}

// OpenURL attempts to open the specified URL in the default web browser.
func OpenURL(url string) error {
	var command string
	var args []string
	switch runtime.GOOS {
	case "linux":
		command = "xdg-open"
		args = []string{url}
	case "darwin":
		command = "open"
		args = []string{url}
	case "windows":
		command = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open URL '%s' using command '%s %v': %w", url, command, args, err)
	}
	return nil
}
