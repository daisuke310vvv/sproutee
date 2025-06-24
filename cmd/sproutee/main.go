package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/daisuke310vvv/sproutee/internal/config"
	"github.com/daisuke310vvv/sproutee/internal/copy"
	"github.com/daisuke310vvv/sproutee/internal/worktree"
	"github.com/spf13/cobra"
)

const (
	osLinux = "linux"
)

var rootCmd = &cobra.Command{
	Use:   "sproutee",
	Short: "A CLI tool for managing Git worktrees efficiently",
	Long: `Sproutee is a CLI tool that automates worktree creation and 
copies specified files to new worktrees based on configuration.

It helps manage multiple branches efficiently by creating worktrees
in .git/sproutee-worktrees/ directory and automatically copying configured files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sproutee - Git Worktree Management Tool")
		fmt.Println("Use 'sproutee --help' for more information.")
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new worktree with file copying",
	Long: `Create a new Git worktree with the specified name. The name will be used 
as both the worktree directory name and the branch name. Files specified in the 
configuration will be automatically copied to the new worktree.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		branch := name

		manager, err := worktree.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Creating worktree '%s' with branch '%s'...\n", name, branch)

		worktreePath, err := manager.CreateWorktree(name, branch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Worktree created successfully at: %s\n", worktreePath)

		fmt.Println("\n📁 Copying configured files...")
		copyReport, err := copy.CopyFilesToWorktree(manager.RepoRoot, worktreePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to copy files: %v\n", err)
		} else {
			copyReport.PrintSummary()
		}

		// Get flags
		openCursor, _ := cmd.Flags().GetBool("cursor")
		openVSCode, _ := cmd.Flags().GetBool("vscode")
		openXcode, _ := cmd.Flags().GetBool("xcode")
		openAndroidStudio, _ := cmd.Flags().GetBool("android-studio")
		customDir, _ := cmd.Flags().GetString("dir")

		// Determine target path for editor
		targetPath := worktreePath
		if customDir != "" {
			if !filepath.IsAbs(customDir) {
				// Relative path: resolve relative to worktree directory
				targetPath = filepath.Join(worktreePath, customDir)
			} else {
				// Absolute path: use as-is
				targetPath = customDir
			}

			// Check if the target path exists
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				fmt.Printf("Warning: Directory '%s' does not exist, using worktree root instead\n", targetPath)
				targetPath = worktreePath
			}
		}

		// Auto-open editor if any flag is set
		if openCursor {
			fmt.Println("\n🚀 Opening Cursor...")
			if customDir != "" {
				fmt.Printf("📁 Target directory: %s\n", targetPath)
			}
			if err := openInEditor(targetPath, "cursor"); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to open Cursor: %v\n", err)
			} else {
				fmt.Println("✅ Cursor opened successfully")
			}
		} else if openVSCode {
			fmt.Println("\n🚀 Opening VS Code...")
			if customDir != "" {
				fmt.Printf("📁 Target directory: %s\n", targetPath)
			}
			if err := openInEditor(targetPath, "vscode"); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to open VS Code: %v\n", err)
			} else {
				fmt.Println("✅ VS Code opened successfully")
			}
		} else if openXcode {
			fmt.Println("\n🚀 Opening Xcode...")
			if customDir != "" {
				fmt.Printf("📁 Target directory: %s\n", targetPath)
			}
			if err := openInEditor(targetPath, "xcode"); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to open Xcode: %v\n", err)
			} else {
				fmt.Println("✅ Xcode opened successfully")
			}
		} else if openAndroidStudio {
			fmt.Println("\n🚀 Opening Android Studio...")
			if customDir != "" {
				fmt.Printf("📁 Target directory: %s\n", targetPath)
			}
			if err := openInEditor(targetPath, "android-studio"); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to open Android Studio: %v\n", err)
			} else {
				fmt.Println("✅ Android Studio opened successfully")
			}
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long:  "Manage Sproutee configuration files and settings.",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  "Create a default sproutee.json configuration file in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get current directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(wd, config.ConfigFileName)

		if err := config.CreateDefaultConfigFile(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration file created: %s\n", configPath)
		fmt.Println("You can now customize the file to specify which files to copy to new worktrees.")
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show configuration",
	Long:  "Display the current configuration settings.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfigFromCurrentDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Current configuration:")
		fmt.Printf("Files to copy: %d\n", len(cfg.CopyFiles))
		for i, file := range cfg.CopyFiles {
			fmt.Printf("  %d. %s\n", i+1, file)
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing worktrees",
	Long:  "Display all existing worktrees created by Sproutee.",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := worktree.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		worktrees, err := manager.ListWorktrees()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(worktrees) == 0 {
			fmt.Println("No worktrees found.")
			return
		}

		fmt.Printf("Found %d worktree(s):\n", len(worktrees))
		for i, wt := range worktrees {
			fmt.Printf("  %d. %s", i+1, wt.Path)
			if wt.Branch != "" {
				fmt.Printf(" (branch: %s)", wt.Branch)
			}
			if wt.Commit != "" {
				fmt.Printf(" [%s]", wt.Commit[:8])
			}
			fmt.Println()
		}
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up worktrees",
	Long:  "Remove unused or orphaned worktrees. Interactive selection with safety checks for uncommitted changes.",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")

		manager, err := worktree.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		worktrees, err := manager.ListWorktrees()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Filter out main worktree (repository root)
		var cleanableWorktrees []worktree.WorktreeInfo
		for _, wt := range worktrees {
			if wt.Path != manager.RepoRoot {
				cleanableWorktrees = append(cleanableWorktrees, wt)
			}
		}

		if len(cleanableWorktrees) == 0 {
			fmt.Println("📁 No additional worktrees found to clean.")
			return
		}

		fmt.Printf("🔍 Found %d worktree(s) to analyze:\n\n", len(cleanableWorktrees))

		// Analyze each worktree
		type worktreeAnalysis struct {
			Info   worktree.WorktreeInfo
			Status *worktree.WorktreeStatus
			Index  int
		}

		var analyses []worktreeAnalysis
		for i, wt := range cleanableWorktrees {
			fmt.Printf("Checking %d. %s...\n", i+1, filepath.Base(wt.Path))

			status, err := manager.CheckWorktreeStatus(wt.Path)
			if err != nil {
				fmt.Printf("   ❌ Error checking status: %v\n", err)
				continue
			}

			analyses = append(analyses, worktreeAnalysis{
				Info:   wt,
				Status: status,
				Index:  i + 1,
			})

			fmt.Printf("   %s\n", status.GetStatusSummary())
			if !status.IsClean() && !force {
				if status.HasStagedChanges || status.HasUnstagedChanges {
					fmt.Printf("   📝 Changed files: %s\n", strings.Join(status.ChangedFiles, ", "))
				}
				if status.HasUntrackedFiles {
					fmt.Printf("   📄 Untracked files: %s\n", strings.Join(status.UntrackedFiles, ", "))
				}
			}
			fmt.Println()
		}

		if len(analyses) == 0 {
			fmt.Println("❌ No worktrees could be analyzed.")
			return
		}

		// Interactive selection
		if !dryRun {
			fmt.Println("💡 Select worktrees to delete:")
			fmt.Println("   - Enter numbers separated by commas (e.g., 1,3,5)")
			fmt.Println("   - Enter 'clean' to delete only clean worktrees")
			fmt.Println("   - Enter 'all' to delete all worktrees")
			fmt.Println("   - Enter 'cancel' to abort")

			if !force {
				fmt.Println("   ⚠️  Worktrees with uncommitted changes will require confirmation")
			}

			fmt.Print("\nYour choice: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "cancel" {
				fmt.Println("❌ Operation cancelled.")
				return
			}

			var selectedIndices []int
			if input == "all" {
				for _, analysis := range analyses {
					selectedIndices = append(selectedIndices, analysis.Index)
				}
			} else if input == "clean" {
				for _, analysis := range analyses {
					if analysis.Status.IsClean() {
						selectedIndices = append(selectedIndices, analysis.Index)
					}
				}
				if len(selectedIndices) == 0 {
					fmt.Println("📁 No clean worktrees found.")
					return
				}
			} else {
				parts := strings.Split(input, ",")
				for _, part := range parts {
					if idx, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
						if idx >= 1 && idx <= len(analyses) {
							selectedIndices = append(selectedIndices, idx)
						}
					}
				}
			}

			if len(selectedIndices) == 0 {
				fmt.Println("❌ No valid worktrees selected.")
				return
			}

			// Process deletions
			fmt.Printf("\n🗑️  Removing %d worktree(s):\n", len(selectedIndices))
			for _, idx := range selectedIndices {
				analysis := analyses[idx-1]
				fmt.Printf("\n🔄 Processing: %s\n", filepath.Base(analysis.Info.Path))

				if !analysis.Status.IsClean() && !force {
					fmt.Printf("⚠️  This worktree has uncommitted changes!\n")
					fmt.Printf("   %s\n", analysis.Status.GetStatusSummary())
					fmt.Print("   Continue with deletion? (y/N): ")

					confirmInput, _ := reader.ReadString('\n')
					if strings.ToLower(strings.TrimSpace(confirmInput)) != "y" {
						fmt.Println("   ⏭️  Skipped.")
						continue
					}
				}

				var removeErr error
				if force || !analysis.Status.IsClean() {
					removeErr = manager.ForceRemoveWorktree(analysis.Info.Path)
				} else {
					removeErr = manager.RemoveWorktree(analysis.Info.Path)
				}

				if removeErr != nil {
					fmt.Printf("   ❌ Failed: %v\n", removeErr)
				} else {
					fmt.Printf("   ✅ Deleted: %s\n", filepath.Base(analysis.Info.Path))
				}
			}
		} else {
			fmt.Println("🔍 Dry run - no worktrees will be deleted:")
			for _, analysis := range analyses {
				status := "would delete"
				if !analysis.Status.IsClean() && !force {
					status = "would require confirmation"
				}
				fmt.Printf("   %d. %s - %s\n", analysis.Index, filepath.Base(analysis.Info.Path), status)
			}
		}
	},
}

// openInEditor opens the specified directory in the chosen editor
func openInEditor(path, editor string) error {
	var cmd *exec.Cmd

	switch editor {
	case "cursor":
		switch runtime.GOOS {
		case "darwin", "windows", osLinux:
			cmd = exec.Command("cursor", path)
		default:
			return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
		}
	case "vscode":
		switch runtime.GOOS {
		case "darwin", "windows", osLinux:
			cmd = exec.Command("code", path)
		default:
			return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
		}
	case "xcode":
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("xed", path)
		default:
			return fmt.Errorf("Xcode is only available on macOS")
		}
	case "android-studio":
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", "-a", "Android Studio", path)
		case "windows":
			cmd = exec.Command("studio", path)
		case osLinux:
			cmd = exec.Command("studio.sh", path)
		default:
			return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
		}
	default:
		return fmt.Errorf("unsupported editor: %s", editor)
	}

	return cmd.Start()
}

func init() {
	createCmd.Flags().Bool("cursor", false, "Automatically open the created worktree in Cursor")
	createCmd.Flags().Bool("vscode", false, "Automatically open the created worktree in VS Code")
	createCmd.Flags().Bool("xcode", false, "Automatically open the created worktree in Xcode (macOS only)")
	createCmd.Flags().Bool("android-studio", false, "Automatically open the created worktree in Android Studio")
	createCmd.Flags().String("dir", "", "Specify directory to open in editor (absolute or relative path)")

	cleanCmd.Flags().Bool("dry-run", false, "Show what would be deleted without actually deleting")
	cleanCmd.Flags().Bool("force", false, "Force deletion without confirmation for worktrees with uncommitted changes")

	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configListCmd)

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(cleanCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
