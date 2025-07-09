# gcp-launch

`gcp-launch` is a command-line utility written in Go that simplifies opening Google Cloud Platform (GCP) console URLs for various services and environments. It supports both a direct command-line interface (CLI) and an interactive Terminal User Interface (TUI).

## Features

*   **Quick Access**: Instantly open GCP console URLs for configured services and environments.
*   **Configurable**: Define your GCP projects, services, and environments in a simple YAML file.
*   **CLI Mode**: Direct command-line usage for scripting and quick launches.
*   **TUI Mode**: Interactive terminal interface for easy navigation and selection.
*   **Autocompletion**: CLI mode supports shell autocompletion for services and environments.
*   **Custom Config Path**: Specify a custom path for your configuration file.

## Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/tom-gray/gcp-launch.git
    cd gcp-launch
    ```

2.  **Build the executable**:
    ```bash
    go build -o gcp-launch .
    ```

3.  **Place the executable in your PATH**:
    Move the `gcp-launch` executable to a directory included in your system's `PATH` (e.g., `/usr/local/bin` or `~/bin`).

    ```bash
    mv gcp-launch /usr/local/bin/
    ```

## Configuration

`gcp-launch` reads its configuration from a YAML file named `.gcp-launch.yaml`. By default, it looks for this file in the same directory as the `gcp-launch` executable.

You can also specify a custom path for the configuration file using the `--config` flag.

### Example `.gcp-launch.yaml`

```yaml
services:
  logging:
    environments:
      myproject-prod:
        project_id: my-prod-project
      myproject-dev:
        project_id: my-dev-project
  cloudrun:
    environments:
      myproject-prod:
        project_id: my-prodk-project
        region: us-central1 # Region is important for Cloud Run URLs
      myproject-dev:
        project_id: my-dev-project
        region: us-central1
  gke:
    environments:
      apps-prod:
        project_id: my-prodk-project
        cluster: my-prod-cluster # Cluster name is important for GKE cluster details
      apps-dev:
        project_id: my-dev-project
        cluster: my-dev-cluster
  spanner:
    environments:
      prod:
        project_id: my-prod-project
      dev:
        project_id: my-dev-project
```

*   `services`: Top-level key containing definitions for different GCP services.
*   `<service_name>`: (e.g., `logging`, `cloudrun`, `gke`) - The name of the GCP service.
*   `environments`: Contains different deployment environments for a service.
*   `<environment_name>`: (e.g., `myproject-prod`, `myproject-dev`) - The name of the environment.
*   `project_id`: The GCP project ID associated with the environment.
*   `region`: (Optional, but recommended for Cloud Run) The GCP region for the service.
*   `cluster`: (Optional, but recommended for GKE) The GKE cluster name.

## Usage

### CLI Mode

Run `gcp-launch` with the service type and environment as arguments.

**Syntax:**

```bash
gcp-launch <service> <environment> [--config <path_to_config>]
```

**Examples:**

1.  **Open Cloud Logging for `myproject-prod`:**
    ```bash
    gcp-launch logging myproject-prod
    ```

2.  **Open Cloud Run for `myproject-dev` in `us-central1`:**
    ```bash
    gcp-launch cloudrun myproject-dev
    ```
    *(Note: The region is automatically used from the config for Cloud Run URLs)*

3.  **Open GKE cluster details for `apps-prod`:**
    ```bash
    gcp-launch gke apps-prod
    ```
    *(Note: The cluster name is automatically used from the config for GKE URLs)*

4.  **Using a custom configuration file:**
    ```bash
    gcp-launch logging myproject-prod --config /path/to/my/custom-config.yaml
    ```

#### Autocompletion

`gcp-launch` supports shell autocompletion. To enable it, you typically need to add a line to your shell's configuration file (e.g., `.bashrc`, `.zshrc`).

For `bash`:
```bash
gcp-launch completion bash >> ~/.bashrc
source ~/.bashrc
```

For `zsh`:
```bash
gcp-launch completion zsh >> ~/.zshrc
source ~/.zshrc
```

*(Refer to Cobra's documentation for more advanced completion setup for other shells if needed.)*

### TUI Mode

Run `gcp-launch` without any arguments to launch the interactive TUI.

```bash
gcp-launch
```

**Navigation:**

*   Use `↑` (up arrow) and `↓` (down arrow) to navigate through the lists.
*   Press `Enter` to select a service or environment.
*   Press `Esc` or `Backspace` to go back to the previous selection.
*   Press `q` or `Ctrl+C` to quit the application.
