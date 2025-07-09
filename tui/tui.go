package tui

import (
	"fmt"
	"sort"
	"strings"

	// "os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tom-gray/gcp-launch/config"
	"github.com/tom-gray/gcp-launch/url"
)

const (
	stateSelectService     = "select_service"
	stateSelectEnvironment = "select_environment"
)

type Model struct {
	cfg               *config.Config
	state             string
	serviceKeys       []string
	serviceCursor     int
	selectedService   string
	environmentKeys   []string
	environmentCursor int
	finalURL          string
	finalError        error
}

func NewModel(cfg *config.Config) Model {
	keys := []string{}
	if cfg != nil && cfg.Services != nil {
		keys = make([]string, 0, len(cfg.Services))
		for key := range cfg.Services {
			keys = append(keys, key)
		}
		sort.Strings(keys)
	}
	return Model{
		cfg:               cfg,
		state:             stateSelectService,
		serviceKeys:       keys,
		serviceCursor:     0,
		selectedService:   "",
		environmentKeys:   nil,
		environmentCursor: 0,
		finalURL:          "",
		finalError:        nil,
	}
}
func (m Model) Init() tea.Cmd { return nil }

// Update handles messages and state transitions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg == nil {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

		switch m.state {
		case stateSelectService:
			// Service selection logic (remains the same)
			switch msg.String() {
			case "up", "k":
				if len(m.serviceKeys) > 0 && m.serviceCursor > 0 {
					m.serviceCursor--
				}
			case "down", "j":
				if len(m.serviceKeys) > 0 && m.serviceCursor < len(m.serviceKeys)-1 {
					m.serviceCursor++
				}
			case "enter":
				if len(m.serviceKeys) > 0 && m.serviceCursor >= 0 && m.serviceCursor < len(m.serviceKeys) {
					m.selectedService = m.serviceKeys[m.serviceCursor]
					envKeys := []string{}
					if serviceConf, ok := m.cfg.Services[m.selectedService]; ok && serviceConf.Environments != nil {
						envKeys = make([]string, 0, len(serviceConf.Environments))
						for k := range serviceConf.Environments {
							envKeys = append(envKeys, k)
						}
						sort.Strings(envKeys)
					}
					m.environmentKeys = envKeys
					m.environmentCursor = 0
					m.state = stateSelectEnvironment
				}
			}

		case stateSelectEnvironment:
			// Environment selection logic
			switch msg.String() {
			case "up", "k":
				if len(m.environmentKeys) > 0 && m.environmentCursor > 0 {
					m.environmentCursor--
				}
			case "down", "j":
				if len(m.environmentKeys) > 0 && m.environmentCursor < len(m.environmentKeys)-1 {
					m.environmentCursor++
				}
			case "enter":
				// --- Handle environment selection ---
				if len(m.environmentKeys) > 0 && m.environmentCursor >= 0 && m.environmentCursor < len(m.environmentKeys) {
					selectedEnv := m.environmentKeys[m.environmentCursor]

					// Safely retrieve the environment configuration
					var envConf config.EnvironmentConfig
					var envOk bool
					if serviceConf, serviceOk := m.cfg.Services[m.selectedService]; serviceOk {
						envConf, envOk = serviceConf.Environments[selectedEnv]
					}
					if !envOk {
						m.finalError = fmt.Errorf("internal error: could not find config for service '%s', environment '%s'", m.selectedService, selectedEnv)
						return m, tea.Quit
					}

					// --- Conditional URL Generation ---
					var serviceURL string
					var genErr error // Only for generic case

					if m.selectedService == "cloudrun" {
						// Specific handling for Cloud Run
						projectID := envConf.ProjectID
						region := envConf.Region
						if projectID == "" || region == "" {
							m.finalError = fmt.Errorf("project_id or region not defined in config for service '%s', environment '%s'", m.selectedService, selectedEnv)
							return m, tea.Quit
						}
						serviceURL = url.GenerateCloudRunURL(projectID, region)

					} else if m.selectedService == "gke" {
						// Specific handling for GKE
						projectID := envConf.ProjectID
						cluster := envConf.Cluster
						if projectID == "" || cluster == "" {
							m.finalError = fmt.Errorf("project_id or cluster not defined in config for service '%s', environment '%s'", m.selectedService, selectedEnv)
							return m, tea.Quit
						}
						serviceURL = url.GenerateGKEURL(projectID, cluster)

					} else {
						// Generic handling for other services
						serviceURL, genErr = url.GenerateServiceURL(m.selectedService, envConf)
						if genErr != nil {
							m.finalError = fmt.Errorf("failed to generate URL: %w", genErr)
							return m, tea.Quit // Quit on generation error
						}
					}

					// --- Attempt to Open URL and Quit ---
					m.finalURL = serviceURL // Store the URL

					openErr := url.OpenURL(serviceURL) // Attempt to open
					if openErr != nil {
						m.finalError = fmt.Errorf("failed to open URL in browser: %w", openErr) // Store open error
					}

					return m, tea.Quit // Quit after attempting generation and opening
				}
			case "esc", "backspace":
				m.state = stateSelectService
				m.selectedService = ""
				m.environmentKeys = nil
				m.environmentCursor = 0
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder
	switch m.state {
	case stateSelectService:
		sb.WriteString("Select a Service Type (Use ↑/↓ arrows, Enter to select, q to quit):\n\n")
		if len(m.serviceKeys) == 0 {
			sb.WriteString("No services defined in the configuration file.\n")
		} else {
			for i, serviceKey := range m.serviceKeys {
				cursorIndicator := "  "
				if m.serviceCursor == i {
					cursorIndicator = "> "
				}
				sb.WriteString(cursorIndicator)
				sb.WriteString(serviceKey)
				sb.WriteString("\n")
			}
		}
	case stateSelectEnvironment:
		sb.WriteString(fmt.Sprintf("Select Environment for '%s' (Use ↑/↓, Enter to open, Esc/Backspace back, q to quit):\n\n", m.selectedService))
		if len(m.environmentKeys) == 0 {
			sb.WriteString(fmt.Sprintf("No environments defined for service '%s'.\n", m.selectedService))
		} else {
			for i, envKey := range m.environmentKeys {
				cursorIndicator := "  "
				if m.environmentCursor == i {
					cursorIndicator = "> "
				}
				sb.WriteString(cursorIndicator)
				sb.WriteString(envKey)
				sb.WriteString("\n")
			}
		}
	default:
		sb.WriteString("Unknown application state.\n")
	}
	sb.WriteString("\n(Press 'q' to quit)\n")
	return sb.String()
}
func (m Model) GetFinalURL() string  { return m.finalURL }
func (m Model) GetFinalError() error { return m.finalError }
