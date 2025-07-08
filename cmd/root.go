package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/tom-gray/gcp-launch/config"
	"github.com/tom-gray/gcp-launch/url"
)

var loadedConfig *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gcp-launch <service> <environment> [context_arg]",
	Short: "Launch GCP service URLs based on configuration.",
	Long: `gcp-launch opens the relevant Google Cloud Platform console URL
for a specified service type and environment based on predefined configuration.

Example: gcp-launch logging development`,
	Args: cobra.RangeArgs(2, 3),
	// MODIFIED: Use the renamed completion function
	ValidArgsFunction: contextualArgCompletion,
	RunE:              executeLaunch,
}

// Execute function remains the same
func Execute(cfg *config.Config) error {
	loadedConfig = cfg
	return rootCmd.Execute()
}

// contextualArgCompletion provides autocompletion suggestions for arguments.
// It suggests service names for the first argument and environment names
// (based on the first argument) for the second argument.
// MODIFIED: Renamed and added logic for len(args) == 1.
func contextualArgCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Ensure config is loaded before attempting completion
	if loadedConfig == nil {
		// Cannot provide completions without config
		// Return an error directive if appropriate, or just no completions
		return nil, cobra.ShellCompDirectiveError // Indicate error state
	}

	switch len(args) {
	case 0:
		// --- Completing the first argument (service) ---
		if loadedConfig.Services == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp // No services in config
		}
		// Extract and sort service keys
		keys := make([]string, 0, len(loadedConfig.Services))
		for k := range loadedConfig.Services {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		// Return service keys, disable file completion
		return keys, cobra.ShellCompDirectiveNoFileComp

	case 1:
		// --- Completing the second argument (environment) ---
		// The first argument (service name) is already provided in args[0]
		serviceName := args[0]

		// Find the service config
		if loadedConfig.Services == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp // No services in config
		}
		serviceConf, serviceExists := loadedConfig.Services[serviceName]

		// If the typed service doesn't exist in config, or has no environments, offer no suggestions
		if !serviceExists || serviceConf.Environments == nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// Extract and sort environment keys for the given service
		envKeys := make([]string, 0, len(serviceConf.Environments))
		for k := range serviceConf.Environments {
			envKeys = append(envKeys, k)
		}
		sort.Strings(envKeys)
		// Return environment keys, disable file completion
		return envKeys, cobra.ShellCompDirectiveNoFileComp

	default:
		// --- Completing third argument (region/cluster) or beyond ---
		// No specific completions provided here, use default behavior (e.g., file completion)
		return nil, cobra.ShellCompDirectiveDefault
	}
}

// executeLaunch function remains the same
func executeLaunch(cmd *cobra.Command, args []string) error {
	service := args[0]
	environment := args[1]
	fmt.Printf("Service Type: %s, Environment: %s\n", service, environment)
	serviceConfig, ok := loadedConfig.Services[service]
	if !ok {
		return fmt.Errorf("service type '%s' not found in configuration", service)
	}
	environmentConfig, ok := serviceConfig.Environments[environment]
	if !ok {
		return fmt.Errorf("environment '%s' not found for service type '%s' in configuration", environment, service)
	}
	if environmentConfig.ProjectID == "" {
		return fmt.Errorf("project_id not defined for service type '%s' in environment '%s'", service, environment)
	}
	var serviceURL string
	var genErr error
	if service == "cloudrun" {
		configRegion := environmentConfig.Region
		if configRegion == "" {
			return fmt.Errorf("region not defined in configuration for service '%s' in environment '%s'", service, environment)
		}
		serviceURL = url.GenerateCloudRunURL(environmentConfig.ProjectID, configRegion)
		genErr = nil
		fmt.Printf("Found project ID: %s, Region: %s. Attempting to open Cloud Run...\n", environmentConfig.ProjectID, configRegion)
	} else if service == "gke" {
		configCluster := environmentConfig.Cluster
		if configCluster == "" {
			return fmt.Errorf("cluster not defined in configuration for service '%s' in environment '%s'", service, environment)
		}
		serviceURL = url.GenerateGKEURL(environmentConfig.ProjectID, configCluster)
		genErr = nil
		fmt.Printf("Found project ID: %s, Cluster: %s. Attempting to open GKE Cluster details...\n", environmentConfig.ProjectID, configCluster)
	} else {
		serviceURL, genErr = url.GenerateServiceURL(service, environmentConfig)
		if genErr != nil {
			return fmt.Errorf("failed to generate URL: %w", genErr)
		}
		fmt.Printf("Found project ID: %s. Attempting to open GCP console for %s...\n", environmentConfig.ProjectID, service)
	}
	openErr := url.OpenURL(serviceURL)
	if openErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to open URL '%s' in browser: %v\n", serviceURL, openErr)
		fmt.Printf("You can manually access the URL here: %s\n", serviceURL)
	} else {
		fmt.Println("Browser should be opening the URL shortly.")
	}
	return nil
}
