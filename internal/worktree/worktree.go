package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	WorktreeDir = ".git/worktree"
)

type Manager struct {
	RepoRoot string
}

func NewManager() (*Manager, error) {
	repoRoot, err := FindGitRepository()
	if err != nil {
		return nil, err
	}
	return &Manager{RepoRoot: repoRoot}, nil
}

func FindGitRepository() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	currentDir, err := filepath.Abs(wd)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	for {
		gitDir := filepath.Join(currentDir, ".git")
		if stat, err := os.Stat(gitDir); err == nil {
			if stat.IsDir() {
				return currentDir, nil
			}
			
			data, err := os.ReadFile(gitDir)
			if err == nil && strings.HasPrefix(string(data), "gitdir: ") {
				return currentDir, nil
			}
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("not inside a Git repository")
}

func generateTimestamp() string {
	return time.Now().Format("20060102_150405")
}

func (m *Manager) GenerateWorktreeDirName(name string) (string, error) {
	timestamp := generateTimestamp()
	return fmt.Sprintf("%s_%s", name, timestamp), nil
}

func (m *Manager) GetWorktreeBasePath() string {
	return filepath.Join(m.RepoRoot, WorktreeDir)
}

func (m *Manager) CreateWorktree(name, branch string) (string, error) {
	dirName, err := m.GenerateWorktreeDirName(name)
	if err != nil {
		return "", fmt.Errorf("failed to generate directory name: %w", err)
	}

	worktreeBasePath := m.GetWorktreeBasePath()
	if err := os.MkdirAll(worktreeBasePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create worktree base directory: %w", err)
	}

	worktreePath := filepath.Join(worktreeBasePath, dirName)

	cmd := exec.Command("git", "worktree", "add", worktreePath, branch)
	cmd.Dir = m.RepoRoot
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create worktree: %w\nOutput: %s", err, string(output))
	}

	return worktreePath, nil
}

func (m *Manager) ListWorktrees() ([]WorktreeInfo, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = m.RepoRoot
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	return parseWorktreeList(string(output))
}

type WorktreeInfo struct {
	Path   string
	Branch string
	Commit string
}

func parseWorktreeList(output string) ([]WorktreeInfo, error) {
	var worktrees []WorktreeInfo
	lines := strings.Split(strings.TrimSpace(output), "\n")
	
	var current WorktreeInfo
	for _, line := range lines {
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = WorktreeInfo{}
			}
			continue
		}
		
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		
		key, value := parts[0], parts[1]
		switch key {
		case "worktree":
			current.Path = value
		case "branch":
			current.Branch = strings.TrimPrefix(value, "refs/heads/")
		case "HEAD":
			current.Commit = value
		}
	}
	
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}
	
	return worktrees, nil
}

func (m *Manager) RemoveWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	cmd.Dir = m.RepoRoot
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %w\nOutput: %s", err, string(output))
	}
	
	return nil
}