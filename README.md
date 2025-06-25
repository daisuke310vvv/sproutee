# Sproutee 🌱

A powerful CLI tool for efficient Git worktree management with automated file copying and multi-editor integration.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](#testing)

## Overview

Sproutee streamlines your Git workflow by automating worktree creation and intelligently copying specified files to new worktrees. Perfect for developers who work with multiple branches simultaneously and need consistent development environments across worktrees.

### Key Features

- 🚀 **Automated Worktree Creation**: Creates Git worktrees with timestamp-based naming
- 📁 **Smart File Copying**: Automatically copies configured files to new worktrees
- 🎯 **Multi-Editor Support**: Launch Cursor, VS Code, Xcode, or Android Studio automatically with custom directory targeting
- ⚙️ **Flexible Configuration**: JSON-based configuration for file management
- 🔧 **Init Script Execution**: Run custom initialization scripts after worktree creation
- 🧹 **Safe Cleanup**: Interactive worktree cleanup with uncommitted change detection
- 🔍 **Status Monitoring**: Track worktree status and changes
- 🛡️ **Cross-Platform**: Works on macOS, Windows, and Linux

## Quick Start

### Installation

#### Homebrew (Recommended)

```bash
# Add the tap
brew tap daisuke310vvv/sproutee

# Install sproutee
brew install sproutee
```

#### Build from Source

```bash
# Clone and build
git clone https://github.com/daisuke310vvv/sproutee.git
cd sproutee
go build -o sproutee cmd/sproutee/main.go

# Move to your PATH (optional)
mv sproutee /usr/local/bin/
```

#### Download Pre-built Binaries

Download the latest release from [GitHub Releases](https://github.com/daisuke310vvv/sproutee/releases).

### Basic Usage

```bash
# Initialize configuration
sproutee config init

# Create a worktree (creates branch with same name)
sproutee create feature-auth

# Create worktree and open in VS Code
sproutee create bugfix-login --vscode

# Create worktree and open specific directory in editor
sproutee create feature-frontend --cursor --dir ./src/frontend

# List all worktrees
sproutee list

# Clean up worktrees interactively
sproutee clean
```

## Commands

### `sproutee create <name>`

Creates a new Git worktree with automatic file copying. The name is used as both the worktree directory name and the branch name.

```bash
# Basic usage - creates worktree and branch both named 'feature-dashboard'
sproutee create feature-dashboard

# With editor integration
sproutee create feature-auth --cursor    # Open in Cursor
sproutee create hotfix-bug --vscode      # Open in VS Code
sproutee create ios-feature --xcode      # Open in Xcode (macOS only)
sproutee create android-fix --android-studio  # Open in Android Studio

# Open specific directory in editor
sproutee create feature-api --cursor --dir ./backend
sproutee create ui-components --vscode --dir ./src/components
```

**Options:**
- `--cursor`: Open worktree in Cursor editor
- `--vscode`: Open worktree in VS Code
- `--xcode`: Open worktree in Xcode (macOS only)
- `--android-studio`: Open worktree in Android Studio
- `--dir <path>`: Specify directory to open in editor (absolute or relative path)

### `sproutee config`

Manage configuration settings.

```bash
sproutee config init    # Create default configuration file
sproutee config list    # Show current configuration
```

### `sproutee list`

Display all existing worktrees with branch and commit information.

```bash
sproutee list
# Output:
# Found 2 worktree(s):
#   1. ~/.sproutee/my-project/feature_20241212_143022 (branch: feature-auth) [a1b2c3d4]
#   2. ~/.sproutee/my-project/bugfix_20241212_144055 (branch: bugfix-login) [e5f6g7h8]
```

### `sproutee clean`

Interactively clean up worktrees with safety checks.

```bash
sproutee clean                    # Interactive cleanup
sproutee clean --dry-run          # Preview what would be deleted
sproutee clean --force            # Skip confirmation for dirty worktrees
```

**Features:**
- Detects uncommitted changes
- Shows file status for each worktree
- Interactive selection (by number, 'clean', or 'all')
- Safety confirmations for worktrees with changes

## Configuration

Sproutee uses a `sproutee.json` configuration file to define which files to copy to new worktrees.

### Configuration File Location

Sproutee searches for `sproutee.json` in the following order:
1. Current directory
2. Parent directories (up to repository root)

### Configuration Options

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `copy_files` | `string[]` | Yes | Array of file paths to copy to new worktrees |
| `init_scripts` | `string[]` | No | Array of commands to execute after worktree creation |

### Configuration Format

```json
{
  "copy_files": [
    ".env",
    ".env.local",
    "docker-compose.yml",
    "package-lock.json",
    "yarn.lock",
    "Makefile",
    ".vscode/settings.json"
  ],
  "init_scripts": ["npm install", "npm run build"]
}
```

### Configuration Examples

**Node.js Project:**
```json
{
  "copy_files": [
    ".env",
    ".env.local",
    "package-lock.json",
    "yarn.lock",
    ".nvmrc"
  ],
  "init_scripts": ["npm install", "npm run build"]
}
```

**Go Project:**
```json
{
  "copy_files": [
    ".env",
    "docker-compose.yml",
    "Makefile",
    ".tool-versions"
  ],
  "init_scripts": ["go mod download"]
}
```

**Python Project:**
```json
{
  "copy_files": [
    ".env",
    "requirements.txt",
    "poetry.lock",
    ".python-version"
  ],
  "init_scripts": ["pip install -r requirements.txt"]
}
```

**Empty Configuration (no file copying or init scripts):**
```json
{
  "copy_files": []
}
```

## Directory Structure

Sproutee organizes worktrees in a clean, predictable structure:

```
your-repo/
├── .git/
│   └── worktrees/                   # Git metadata (managed by Git)
│       ├── feature_20241212_143022/
│       └── bugfix_20241212_144055/
├── sproutee.json                    # Configuration file
└── ...                             # Your project files

~/.sproutee/                         # Sproutee home directory
└── your-repo/                       # Project-specific worktrees
    ├── feature_20241212_143022/     # Actual worktree code
    └── bugfix_20241212_144055/      # Actual worktree code
```

## Editor Integration

Sproutee supports automatic editor launching for popular development environments:

| Editor | Command | Platforms | Notes |
|--------|---------|-----------|-------|
| **Cursor** | `cursor` | macOS, Windows, Linux | AI-powered editor |
| **VS Code** | `code` | macOS, Windows, Linux | Microsoft Visual Studio Code |
| **Xcode** | `xed` | macOS only | Apple's IDE for iOS/macOS development |
| **Android Studio** | `studio` / `open -a` | All platforms | Google's Android IDE |

**Requirements:**
- Respective command-line tools must be installed
- Editors must be accessible from PATH

### Directory Targeting

The `--dir` option allows you to specify which directory to open in your editor:

```bash
# Open specific subdirectory (relative path)
sproutee create feature-backend --cursor --dir ./backend

# Open specific subdirectory (absolute path)
sproutee create feature-docs --vscode --dir /path/to/documentation

# Without --dir, opens the worktree root (default behavior)
sproutee create feature --cursor
```

**Behavior:**
- **Relative paths**: Resolved relative to the worktree directory
- **Absolute paths**: Used as-is
- **Non-existent paths**: Falls back to worktree root with a warning
- **Without --dir**: Opens the worktree root directory (default)

## Examples

### Typical Workflow

```bash
# 1. Set up project configuration
cd my-project
sproutee config init

# 2. Edit sproutee.json to include necessary files and init scripts
echo '{
  "copy_files": [
    ".env",
    "docker-compose.yml",
    "package-lock.json"
  ],
  "init_scripts": ["npm install", "npm run build"]
}' > sproutee.json

# 3. Create feature worktree
sproutee create feature-user-auth --vscode

# 4. Create frontend-specific worktree
sproutee create feature-ui --cursor --dir ./frontend

# 5. Work on feature...

# 6. Create another worktree for hotfix
sproutee create hotfix-critical-bug --cursor

# 6. View all worktrees
sproutee list

# 7. Clean up when done
sproutee clean
```

### Advanced Usage

```bash
# Create worktree (automatically creates branch with same name)
sproutee create feature-new-api

# Open specific subdirectory in editor
sproutee create backend-refactor --vscode --dir ./api
sproutee create mobile-app --android-studio --dir ./mobile

# Use absolute paths
sproutee create docs-update --cursor --dir /path/to/docs

# Clean up only clean worktrees
sproutee clean
# Then select: clean

# Force cleanup without confirmations
sproutee clean --force
```

## Development

### Prerequisites

- Go 1.21 or higher
- Git 2.5 or higher

### Building from Source

```bash
git clone https://github.com/daisuke310vvv/sproutee.git
cd sproutee
go mod download
go build -o sproutee cmd/sproutee/main.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/config
go test ./internal/copy
go test ./internal/worktree
```

### Project Structure

```
sproutee/
├── cmd/sproutee/           # Main application entry point
│   └── main.go
├── internal/               # Internal packages
│   ├── config/            # Configuration management
│   ├── copy/              # File copying operations
│   └── worktree/          # Git worktree operations
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md              # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go conventions and idioms
- Add tests for new functionality
- Update documentation as needed
- Use meaningful commit messages

## Troubleshooting

### Common Issues

**Error: "not inside a Git repository"**
- Ensure you're running Sproutee from within a Git repository
- Check that `.git` directory exists in current or parent directories

**Error: "configuration file 'sproutee.json' not found"**
- Run `sproutee config init` to create a default configuration
- Ensure `sproutee.json` exists in current directory or parent directories

**Editor fails to open**
- Verify the editor's command-line tool is installed and in PATH
- For VS Code: Install "code" command via Command Palette
- For Cursor: Ensure "cursor" command is available
- For Xcode: "xed" should be available with Xcode installation

**Permission denied errors**
- Ensure you have write permissions in the repository directory
- Check file ownership and permissions for `.git` directory

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Inspired by Git's powerful worktree functionality
- Thanks to the open-source community for tools and inspiration

---

**Made with ❤️ by [daisuke310vvv](https://github.com/daisuke310vvv)**