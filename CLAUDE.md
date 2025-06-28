# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Run
- `go build .` - Build the CLI application
- `go run .` - Run the CLI application directly
- `go run . --help` - Show available commands and flags
- `go run . --config /path/to/config.yaml` - Use specific configuration file

### Configuration
- `~/.config/cli/config.yaml` - Default configuration file location
- `--config` flag - Override default configuration file path
- `sample.yaml` - Example configuration file showing structure

### Application Commands
- `go run . build <app-name>` - Build application using Docker
- `go run . build <app-name> --verbose` - Build with full Docker output
- `go run . deploy app <app-name>` - Deploy application using Helm
- `go run . deploy app <app-name> --verbose` - Deploy with full Helm output
- `go run . deploy dependencies` - Deploy dependencies using Helm
- `go run . deploy dependencies --optional` - Include optional dependencies
- `go run . deploy dependencies --verbose` - Deploy with full Helm output

### Cluster Commands
- `go run . cluster create` - Create a new cluster
- `go run . cluster delete` - Delete cluster (with confirmation)
- `go run . cluster start` - Start an existing cluster
- `go run . cluster stop` - Stop a running cluster

### Development
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download dependencies
- `devbox shell` - Enter development environment (uses Nix/devbox for dependency management)

### Testing
- `go test ./...` - Run all tests (when tests are added)
- `go test -v ./...` - Run tests with verbose output

## Architecture

This is a Go CLI application built using the Cobra framework with clear separation of concerns:

### Project Structure
```
go-cli/
├── main.go                    # Entry point - calls cmd.Execute()
├── cmd/
│   ├── root.go               # Root command + global config (Viper)
│   ├── build/build.go        # Build command UI layer
│   ├── deploy/deploy.go      # Deploy command UI layer
│   └── cluster/cluster.go    # Cluster command UI layer
├── internal/
│   ├── build/build.go        # Build business logic (Docker integration)
│   ├── deploy/deploy.go      # Deploy business logic (Helm integration)
│   └── cluster/cluster.go    # Cluster business logic
├── sample.yaml               # Example configuration file
└── CLAUDE.md                # This file
```

### Key Components
- **Cobra**: CLI framework for command structure and argument parsing
- **Viper**: Configuration management with YAML support
- **Spinner**: User feedback during long-running operations
- **Docker**: Container image building
- **Helm**: Kubernetes application deployment

### Architecture Patterns
- **Separation of Concerns**: `cmd/` handles UI/UX, `internal/` handles business logic
- **Configuration-Driven**: All build and deploy parameters read from YAML config
- **Command Registration**: Each subcommand uses `GetCommand()` pattern for registration
- **Error Handling**: Consistent error reporting across all commands
- **Verbose Control**: Global `--verbose` flags for detailed output control

### Configuration Structure
Applications are configured under the `apps` key, with each app having:
- `project_path`: Path to application source code
- `build`: Docker build configuration (image_name, dockerfile, context, build_args)
- `deploy`: Helm deployment configuration (chart_path, values_file, namespace)

Dependencies are configured under the `dependencies` key with Helm chart details.

### Tool Integration
- **Docker**: Real `docker build` commands for application building
- **Helm**: Real `helm upgrade --install` commands for deployments
- **Spinners**: Visual feedback during tool execution (hidden when `--verbose` is used)

### Adding New Commands
New subcommands should:
1. Create UI layer in `cmd/[command]/[command].go`
2. Create business logic in `internal/[command]/[command].go`
3. Use `GetCommand()` pattern for registration
4. Follow consistent error handling and verbose flag patterns
5. Register in `cmd/root.go` imports