package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/daisuke310vvv/sproutee/internal/config"
)

type Result struct {
	SourcePath string
	TargetPath string
	Success    bool
	Error      error
}

type Report struct {
	Results      []Result
	TotalFiles   int
	SuccessCount int
	FailureCount int
}

func (r *Report) AddResult(result Result) {
	r.Results = append(r.Results, result)
	r.TotalFiles++
	if result.Success {
		r.SuccessCount++
	} else {
		r.FailureCount++
	}
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func File(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	targetFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	sourceInfo, err := sourceFile.Stat()
	if err == nil {
		_ = targetFile.Chmod(sourceInfo.Mode())
	}

	return nil
}

func FileWithStructure(srcRoot, targetRoot, relativePath string) error {
	srcPath := filepath.Join(srcRoot, relativePath)
	dstPath := filepath.Join(targetRoot, relativePath)

	if !FileExists(srcPath) {
		return fmt.Errorf("source file does not exist: %s", srcPath)
	}

	return File(srcPath, dstPath)
}

func FilesFromConfig(srcRoot, targetRoot string, cfg *config.Config) *Report {
	report := &Report{}

	for _, filePath := range cfg.CopyFiles {
		result := Result{
			SourcePath: filepath.Join(srcRoot, filePath),
			TargetPath: filepath.Join(targetRoot, filePath),
		}

		if !FileExists(result.SourcePath) {
			result.Success = false
			result.Error = fmt.Errorf("source file does not exist: %s", result.SourcePath)
		} else {
			err := FileWithStructure(srcRoot, targetRoot, filePath)
			if err != nil {
				result.Success = false
				result.Error = err
			} else {
				result.Success = true
			}
		}

		report.AddResult(result)
	}

	return report
}

func FilesToWorktree(sourceRepoRoot, worktreePath string) (*Report, error) {
	cfg, err := config.LoadConfigFromCurrentDir()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return FilesFromConfig(sourceRepoRoot, worktreePath, cfg), nil
}

func (r *Report) PrintSummary() {
	if r.TotalFiles == 0 {
		fmt.Println("📁 No files configured for copying.")
		return
	}

	fmt.Printf("📁 File Copy Summary:\n")
	fmt.Printf("   Total files: %d\n", r.TotalFiles)
	fmt.Printf("   ✅ Successful: %d\n", r.SuccessCount)

	if r.FailureCount > 0 {
		fmt.Printf("   ❌ Failed: %d\n", r.FailureCount)
		fmt.Println("\n📋 Failed copies:")
		for _, result := range r.Results {
			if !result.Success {
				fmt.Printf("   • %s → %s\n", result.SourcePath, result.TargetPath)
				fmt.Printf("     Error: %v\n", result.Error)
			}
		}
	}

	if r.SuccessCount > 0 {
		fmt.Println("\n📋 Successfully copied files:")
		for _, result := range r.Results {
			if result.Success {
				relativeTarget := strings.TrimPrefix(result.TargetPath, result.TargetPath[:strings.LastIndex(result.TargetPath, "/")+1])
				fmt.Printf("   • %s\n", relativeTarget)
			}
		}
	}
}
