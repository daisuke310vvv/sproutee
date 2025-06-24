package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGenerateTimestamp(t *testing.T) {
	timestamp1 := generateTimestamp()
	time.Sleep(time.Second)
	timestamp2 := generateTimestamp()

	if timestamp1 == timestamp2 {
		t.Error("generateTimestamp() should generate different timestamps")
	}

	if len(timestamp1) != 15 { // 20060102_150405 format
		t.Errorf("Timestamp format incorrect: %s", timestamp1)
	}
}

func TestFindGitRepository(t *testing.T) {
	tempDir := t.TempDir()

	gitDir := filepath.Join(tempDir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	repoRoot, err := FindGitRepository()
	if err != nil {
		t.Fatalf("FindGitRepository() error = %v", err)
	}

	absRepoRoot, err := filepath.EvalSymlinks(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	absTempDir, err := filepath.EvalSymlinks(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	if absRepoRoot != absTempDir {
		t.Errorf("FindGitRepository() = %v, want %v", absRepoRoot, absTempDir)
	}

	nonGitDir := t.TempDir()
	if err := os.Chdir(nonGitDir); err != nil {
		t.Fatal(err)
	}

	_, err = FindGitRepository()
	if err == nil {
		t.Error("FindGitRepository() should return error for non-git directory")
	}
}

func TestManagerGenerateWorktreeDirName(t *testing.T) {
	manager := &Manager{RepoRoot: "/test"}

	name := "feature-123"
	dirName, err := manager.GenerateWorktreeDirName(name)
	if err != nil {
		t.Fatalf("GenerateWorktreeDirName() error = %v", err)
	}

	if !strings.HasPrefix(dirName, name+"_") {
		t.Errorf("Directory name should start with %s_, got %s", name, dirName)
	}

	if len(dirName) <= len(name)+1 {
		t.Error("Directory name should include timestamp")
	}
}

func TestManagerGetWorktreeBasePath(t *testing.T) {
	repoRoot := "/test/repo"
	manager := &Manager{RepoRoot: repoRoot}

	basePath := manager.GetWorktreeBasePath()
	expected := filepath.Join(repoRoot, WorktreeDir)

	if basePath != expected {
		t.Errorf("GetWorktreeBasePath() = %v, want %v", basePath, expected)
	}
}

func TestParseWorktreeList(t *testing.T) {
	output := `worktree /path/to/main
HEAD 1234567890abcdef
branch refs/heads/main

worktree /path/to/feature
HEAD abcdef1234567890
branch refs/heads/feature-branch

worktree /path/to/detached
HEAD fedcba0987654321
`

	worktrees, err := parseWorktreeList(output)
	if err != nil {
		t.Fatalf("parseWorktreeList() error = %v", err)
	}

	if len(worktrees) != 3 {
		t.Errorf("Expected 3 worktrees, got %d", len(worktrees))
	}

	if worktrees[0].Path != "/path/to/main" {
		t.Errorf("First worktree path = %s, want /path/to/main", worktrees[0].Path)
	}

	if worktrees[0].Branch != "main" {
		t.Errorf("First worktree branch = %s, want main", worktrees[0].Branch)
	}

	if worktrees[0].Commit != "1234567890abcdef" {
		t.Errorf("First worktree commit = %s, want 1234567890abcdef", worktrees[0].Commit)
	}

	if worktrees[1].Branch != "feature-branch" {
		t.Errorf("Second worktree branch = %s, want feature-branch", worktrees[1].Branch)
	}

	if worktrees[2].Branch != "" {
		t.Errorf("Third worktree should have empty branch for detached HEAD, got %s", worktrees[2].Branch)
	}
}

func TestParseWorktreeListEmpty(t *testing.T) {
	output := ""
	worktrees, err := parseWorktreeList(output)
	if err != nil {
		t.Fatalf("parseWorktreeList() error = %v", err)
	}

	if len(worktrees) != 0 {
		t.Errorf("Expected 0 worktrees for empty output, got %d", len(worktrees))
	}
}

func TestContains(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	if !contains(slice, "banana") {
		t.Error("contains() should return true for existing item")
	}

	if contains(slice, "grape") {
		t.Error("contains() should return false for non-existing item")
	}

	if contains([]string{}, "any") {
		t.Error("contains() should return false for empty slice")
	}
}

func TestStatusIsClean(t *testing.T) {
	cleanStatus := &Status{
		HasUnstagedChanges: false,
		HasStagedChanges:   false,
		HasUntrackedFiles:  false,
	}

	if !cleanStatus.IsClean() {
		t.Error("Clean status should return true for IsClean()")
	}

	dirtyStatus := &Status{
		HasUnstagedChanges: true,
		HasStagedChanges:   false,
		HasUntrackedFiles:  false,
	}

	if dirtyStatus.IsClean() {
		t.Error("Dirty status should return false for IsClean()")
	}
}

func TestStatusGetStatusSummary(t *testing.T) {
	cleanStatus := &Status{
		HasUnstagedChanges: false,
		HasStagedChanges:   false,
		HasUntrackedFiles:  false,
	}

	summary := cleanStatus.GetStatusSummary()
	if !strings.Contains(summary, "Clean") {
		t.Errorf("Clean status summary should contain 'Clean', got: %s", summary)
	}

	dirtyStatus := &Status{
		HasUnstagedChanges: true,
		HasStagedChanges:   true,
		HasUntrackedFiles:  true,
		UntrackedFiles:     []string{"file1.txt", "file2.txt"},
	}

	summary = dirtyStatus.GetStatusSummary()
	if !strings.Contains(summary, "staged changes") {
		t.Errorf("Dirty status should mention staged changes, got: %s", summary)
	}
	if !strings.Contains(summary, "unstaged changes") {
		t.Errorf("Dirty status should mention unstaged changes, got: %s", summary)
	}
	if !strings.Contains(summary, "untracked files") {
		t.Errorf("Dirty status should mention untracked files, got: %s", summary)
	}
}

func TestNewManagerWithoutGitRepo(t *testing.T) {
	tempDir := t.TempDir()

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	_, err = NewManager()
	if err == nil {
		t.Error("NewManager() should return error when not in git repository")
	}
}
