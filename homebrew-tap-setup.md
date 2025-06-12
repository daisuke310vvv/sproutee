# Homebrew Tap Setup Guide

This guide explains how to set up the Homebrew tap for Sproutee.

## Prerequisites

1. GoReleaser must be installed
2. You need write access to GitHub repositories
3. A GitHub token with repo permissions

## Setup Steps

### 1. Create Homebrew Tap Repository

Create a new GitHub repository named `homebrew-sproutee`:

```bash
# Create the repository on GitHub
gh repo create daisuke310vvv/homebrew-sproutee --public --description "Homebrew tap for Sproutee CLI tool"

# Clone it locally
git clone https://github.com/daisuke310vvv/homebrew-sproutee.git
cd homebrew-sproutee

# Create initial README
echo "# Homebrew Tap for Sproutee

A powerful CLI tool for efficient Git worktree management.

## Installation

\`\`\`bash
brew tap daisuke310vvv/sproutee
brew install sproutee
\`\`\`

## Usage

\`\`\`bash
sproutee --help
\`\`\`
" > README.md

git add README.md
git commit -m "Initial commit"
git push origin main
```

### 2. Configure GoReleaser

The `.goreleaser.yml` file has been updated to include Homebrew formula generation.

### 3. Set up GitHub Token

Create a GitHub token with `repo` permissions and add it to your repository secrets:

1. Go to GitHub Settings > Developer settings > Personal access tokens
2. Generate a new token with `repo` scope
3. Add it as `HOMEBREW_TAP_GITHUB_TOKEN` in your repository secrets

### 4. Test Release Process

1. Create a test tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. GoReleaser will automatically:
   - Build binaries for multiple platforms
   - Create GitHub release
   - Update the Homebrew formula in your tap repository

### 5. Installation Instructions

Once the tap is set up, users can install Sproutee with:

```bash
brew tap daisuke310vvv/sproutee
brew install sproutee
```

## Maintenance

- Formula will be automatically updated on each release
- Monitor the tap repository for any issues
- Test installations periodically

## Alternative: Manual Formula

If automatic updates don't work, you can manually maintain the formula in the `Formula/sproutee.rb` file in your tap repository.

## Notes

- The tap repository must be named `homebrew-<tool-name>`
- Formula file must be in `Formula/` directory
- GoReleaser handles SHA256 checksum calculation automatically