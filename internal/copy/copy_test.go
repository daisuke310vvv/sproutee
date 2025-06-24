package copy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/daisuke310vvv/sproutee/internal/config"
)

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()

	existingFile := filepath.Join(tempDir, "existing.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	if !FileExists(existingFile) {
		t.Error("FileExists() should return true for existing file")
	}

	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	if FileExists(nonExistentFile) {
		t.Error("FileExists() should return false for non-existent file")
	}
}

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	srcFile := filepath.Join(tempDir, "source.txt")
	testContent := "test file content"
	if err := os.WriteFile(srcFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tempDir, "subdir")
	dstFile := filepath.Join(dstDir, "destination.txt")

	err := CopyFile(srcFile, dstFile)
	if err != nil {
		t.Fatalf("CopyFile() error = %v", err)
	}

	if !FileExists(dstFile) {
		t.Error("Destination file was not created")
	}

	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != testContent {
		t.Errorf("File content = %s, want %s", string(content), testContent)
	}
}

func TestCopyFileWithStructure(t *testing.T) {
	tempDir := t.TempDir()

	srcRoot := filepath.Join(tempDir, "src")
	targetRoot := filepath.Join(tempDir, "target")

	if err := os.MkdirAll(filepath.Join(srcRoot, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}

	relativePath := filepath.Join("subdir", "test.txt")
	srcFile := filepath.Join(srcRoot, relativePath)
	testContent := "structured file content"

	if err := os.WriteFile(srcFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	err := CopyFileWithStructure(srcRoot, targetRoot, relativePath)
	if err != nil {
		t.Fatalf("CopyFileWithStructure() error = %v", err)
	}

	targetFile := filepath.Join(targetRoot, relativePath)
	if !FileExists(targetFile) {
		t.Error("Target file was not created")
	}

	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != testContent {
		t.Errorf("File content = %s, want %s", string(content), testContent)
	}

	nonExistentPath := "nonexistent/file.txt"
	err = CopyFileWithStructure(srcRoot, targetRoot, nonExistentPath)
	if err == nil {
		t.Error("CopyFileWithStructure() should return error for non-existent source")
	}
}

func TestCopyFilesFromConfig(t *testing.T) {
	tempDir := t.TempDir()

	srcRoot := filepath.Join(tempDir, "src")
	targetRoot := filepath.Join(tempDir, "target")

	if err := os.MkdirAll(srcRoot, 0755); err != nil {
		t.Fatal(err)
	}

	existingFile := ".env"
	nonExistentFile := ".nonexistent"

	if err := os.WriteFile(filepath.Join(srcRoot, existingFile), []byte("TEST=value"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		CopyFiles: []string{existingFile, nonExistentFile},
	}

	report := CopyFilesFromConfig(srcRoot, targetRoot, cfg)

	if report.TotalFiles != 2 {
		t.Errorf("TotalFiles = %d, want 2", report.TotalFiles)
	}

	if report.SuccessCount != 1 {
		t.Errorf("SuccessCount = %d, want 1", report.SuccessCount)
	}

	if report.FailureCount != 1 {
		t.Errorf("FailureCount = %d, want 1", report.FailureCount)
	}

	if !FileExists(filepath.Join(targetRoot, existingFile)) {
		t.Error("Existing file was not copied")
	}

	if FileExists(filepath.Join(targetRoot, nonExistentFile)) {
		t.Error("Non-existent file should not have been copied")
	}
}

func TestCopyReport_AddResult(t *testing.T) {
	report := &CopyReport{}

	successResult := CopyResult{
		SourcePath: "/src/file1.txt",
		TargetPath: "/target/file1.txt",
		Success:    true,
		Error:      nil,
	}

	failureResult := CopyResult{
		SourcePath: "/src/file2.txt",
		TargetPath: "/target/file2.txt",
		Success:    false,
		Error:      os.ErrNotExist,
	}

	report.AddResult(successResult)
	report.AddResult(failureResult)

	if report.TotalFiles != 2 {
		t.Errorf("TotalFiles = %d, want 2", report.TotalFiles)
	}

	if report.SuccessCount != 1 {
		t.Errorf("SuccessCount = %d, want 1", report.SuccessCount)
	}

	if report.FailureCount != 1 {
		t.Errorf("FailureCount = %d, want 1", report.FailureCount)
	}

	if len(report.Results) != 2 {
		t.Errorf("Results length = %d, want 2", len(report.Results))
	}
}
