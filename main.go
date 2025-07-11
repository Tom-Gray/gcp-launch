package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tom-gray/gcp-launch/cmd"
	"github.com/tom-gray/gcp-launch/config"
	"github.com/tom-gray/gcp-launch/tui"
)

func main() {
	var configPath string
	var debugMode bool

	// Parse command line arguments for config and debug flags
	// This needs to be done before cobra.Command.Execute() is called
	for i, arg := range os.Args {
		if arg == "--config" && i+1 < len(os.Args) {
			configPath = os.Args[i+1]
			// Remove the flag and its value from os.Args so Cobra doesn't try to parse it
			os.Args = append(os.Args[:i], os.Args[i+2:]...)
			break
		}
	}

	// Check for debug flag (but don't remove it since cobra will handle it)
	for _, arg := range os.Args {
		if arg == "--debug" {
			debugMode = true
			break
		}
	}

	// Check if we should run TUI mode (no service/environment arguments provided)
	// Count non-flag arguments to determine if service and environment were provided
	nonFlagArgs := 0
	for _, arg := range os.Args[1:] { // Skip program name
		if !strings.HasPrefix(arg, "--") {
			nonFlagArgs++
		}
	}

	if debugMode {
		fmt.Printf("[DEBUG] Attempting to load configuration from: %s\n", func() string {
			if configPath == "" {
				return "default location"
			}
			return configPath
		}())
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if nonFlagArgs == 0 {
		// --- TUI Mode ---
		if debugMode {
			fmt.Println("[DEBUG] No arguments provided, launching TUI...")
		}
		initialModel := tui.NewModel(cfg)
		p := tea.NewProgram(initialModel, tea.WithAltScreen())
		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
		if fm, ok := finalModel.(tui.Model); ok {
			finalErr := fm.GetFinalError()
			finalURL := fm.GetFinalURL()
			if finalErr != nil {
				if strings.Contains(finalErr.Error(), "failed to open URL") && finalURL != "" {
					fmt.Fprintf(os.Stderr, "Warning: %v\n", finalErr)
					fmt.Printf("You can manually access the URL here: %s\n", finalURL)
				} else {
					fmt.Fprintf(os.Stderr, "Error during TUI operation: %v\n", finalErr)
				}
			} else if finalURL != "" {
				fmt.Println("Launching:", finalURL)
			} else {
				fmt.Println("TUI finished.")
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: Could not read final TUI state.\n")
			os.Exit(1)
		}
	} else {
		// --- CLI Mode ---

		// Pass the loaded config to the command execution context
		if err := cmd.Execute(cfg); err != nil {
			// Cobra RunE errors are caught here
			// Cobra automatically prints the error, but we might add extra context or just exit
			// fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err) // Optional: Cobra usually prints errors
			os.Exit(1) // Exit on command execution error
		}
		// Command executed successfully
	}
}
