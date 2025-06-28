# Go CLI

A powerful command-line interface for managing application builds and deployments using Docker and Helm.

## Features

- **Application Building**: Build Docker images with configuration-driven parameters
- **Application Deployment**: Deploy applications to Kubernetes using Helm charts
- **Dependency Management**: Deploy and manage Helm chart dependencies
- **Cluster Operations**: Basic cluster lifecycle management
- **Configuration Management**: YAML-based configuration with Viper
- **Verbose Output Control**: Toggle detailed output for all operations

## Quick Start

### Prerequisites

- Go 1.24 or later
- Docker (for building applications)
- Helm (for deployments)
- Kubernetes cluster access (for deployments)

### Installation

```bash
git clone <repository-url>
cd go-cli
go mod download
go build .
```

### Configuration

Create a configuration file at `~/.config/cli/config.yaml` or use the `--config` flag:

```yaml
apps:
  api:
    project_path: "/path/to/api"
    build:
      image_name: api:latest
      dockerfile: Dockerfile
      context: .
      build_args:
        - "PYTHON_VERSION=3.12"
    deploy:
      chart_path: ./helm/chart
      values_file: ./helm/values.yaml
      namespace: application

dependencies:
  redis:
    chart_name: bitnami/redis
    values_file: ./configs/redis-values.yaml
    version: 19.0.0
    namespace: database
  postgresql:
    chart_name: bitnami/postgresql
    values_file: ./configs/postgres-values.yaml
    version: 14.0.0
    namespace: database
    optional: true
```

## Usage

### Building Applications

```bash
# Build an application
./go-cli build api

# Build with verbose Docker output
./go-cli build api --verbose
```

### Deploying Applications

```bash
# Deploy an application
./go-cli deploy app api

# Deploy with verbose Helm output
./go-cli deploy app api --verbose
```

### Managing Dependencies

```bash
# Deploy all dependencies
./go-cli deploy dependencies

# Deploy including optional dependencies
./go-cli deploy dependencies --optional

# Deploy with verbose output
./go-cli deploy dependencies --verbose
```

### Cluster Operations

```bash
# Create a cluster
./go-cli cluster create

# Delete a cluster (with confirmation)
./go-cli cluster delete

# Start/stop cluster
./go-cli cluster start
./go-cli cluster stop
```

### Global Options

```bash
# Use custom configuration file
./go-cli --config /path/to/config.yaml build api

# Show help for any command
./go-cli --help
./go-cli build --help
```

## Architecture

The CLI follows a clean architecture with separation of concerns:

- **`cmd/`**: Command-line interface and user interaction
- **`internal/`**: Business logic and tool integrations
- **Configuration**: YAML-based with Viper for management
- **Tools**: Real integration with Docker and Helm

### Project Structure

```
go-cli/
├── main.go                    # Application entry point
├── cmd/
│   ├── root.go               # Root command and global configuration
│   ├── build/                # Build command implementation
│   ├── deploy/               # Deploy command implementation
│   └── cluster/              # Cluster command implementation
├── internal/
│   ├── build/                # Docker build logic
│   ├── deploy/               # Helm deployment logic
│   └── cluster/              # Cluster management logic
├── sample.yaml               # Example configuration
├── CLAUDE.md                 # Development guidance
└── README.md                 # This file
```

## Configuration Reference

### Application Configuration

Each application under the `apps` key supports:

- **`project_path`**: Path to the application source code
- **`build`**: Docker build configuration
  - `image_name`: Docker image name and tag
  - `dockerfile`: Dockerfile path (relative to context)
  - `context`: Build context path
  - `build_args`: List of build arguments (optional)
- **`deploy`**: Helm deployment configuration
  - `chart_path`: Path to Helm chart
  - `values_file`: Path to values file
  - `namespace`: Kubernetes namespace (required)

### Dependencies Configuration

Dependencies under the `dependencies` key support:

- **`chart_name`**: Helm chart name (e.g., `bitnami/redis`)
- **`values_file`**: Path to values file
- **`version`**: Chart version
- **`namespace`**: Kubernetes namespace
- **`optional`**: Whether dependency is optional (default: false)

## Development

### Prerequisites

- Go 1.24+
- [devbox](https://www.jetpack.io/devbox) (optional, for consistent development environment)

### Setup

```bash
# Using devbox (recommended)
devbox shell

# Or standard Go development
go mod download
```

### Building

```bash
go build .
```

### Testing

```bash
go test ./...
go test -v ./...  # verbose output
```

### Adding New Commands

1. Create UI layer in `cmd/[command]/[command].go`
2. Create business logic in `internal/[command]/[command].go`
3. Use `GetCommand()` pattern for registration
4. Follow existing patterns for error handling and verbose flags
5. Register in `cmd/root.go` imports

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the existing architecture
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License. See the LICENSE file for details.