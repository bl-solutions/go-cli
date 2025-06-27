# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Run
- `go build .` - Build the CLI application
- `go run .` - Run the CLI application directly
- `go run . --help` - Show available commands and flags
- `go run . cluster` - Execute the cluster subcommand

### Development
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download dependencies
- `devbox shell` - Enter development environment (uses Nix/devbox for dependency management)

### Testing
- `go test ./...` - Run all tests (when tests are added)
- `go test -v ./...` - Run tests with verbose output

## Architecture

This is a Go CLI application built using the Cobra framework. The architecture follows Cobra's standard patterns:

- **Entry point**: `main.go` calls `cmd.Execute()` to start the CLI
- **Root command**: `cmd/root.go` defines the base `go-cli` command and global configuration
- **Subcommands**: Located in `cmd/` subdirectories (e.g., `cmd/cluster/`)
- **Command registration**: Each subcommand registers itself with the root command in its `init()` function

### Key Components
- Uses `github.com/spf13/cobra` for CLI framework
- Module name: `go-cli`
- Go version: 1.24
- Development environment managed by devbox with cobra-cli tools

### Adding New Commands
New subcommands should:
1. Create a new package under `cmd/`
2. Define the command using `cobra.Command`
3. Register with `cmd.RootCmd.AddCommand()` in the `init()` function
4. Import the package in `cmd/root.go` or ensure it's loaded