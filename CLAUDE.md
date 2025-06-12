# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sproutee is a CLI tool for managing Git worktrees efficiently. It automates worktree creation and copies specified files to new worktrees based on configuration.

## Development Setup

This project is planned to use Go 1.24.4 with asdf for version management:

```bash
# Setup Go environment
asdf plugin add golang https://github.com/asdf-community/asdf-golang.git
asdf install golang 1.24.4
asdf global golang 1.24.4

# Initialize project (when ready)
go mod init github.com/daisuke310vvv/sproutee
```

## Architecture

The project follows standard Go project structure:
- `cmd/sproutee/main.go` - Entry point
- `internal/` - Internal packages (config, worktree, copy)
- `pkg/` - External packages

Key components:
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Configuration**: JSON files using standard library
- **Worktree Management**: Git worktree operations in `.git/sproutee-worktrees/` directory
- **File Operations**: Automated file copying based on configuration

## Configuration

Uses `sproutee.json` configuration file format:
```json
{
  "copy_files": [
    ".env",
    ".env.local",
    "docker-compose.yml"
  ]
}
```

## Planned Commands

- `sproutee create <name> [branch]` - Create worktree with file copying
- `sproutee config init` - Initialize configuration file
- `sproutee config list` - Show configuration
- `sproutee list` - List existing worktrees
- `sproutee clean` - Clean up worktrees

## Testing

When implementing tests:
- Unit tests for each internal package
- Integration tests for CLI commands
- Test with various Git repository states
- Error handling validation

## Build & Release

Planned to use GoReleaser for multi-platform builds and Homebrew distribution.