package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/daisuke310vvv/sproutee/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sproutee",
	Short: "A CLI tool for managing Git worktrees efficiently",
	Long: `Sproutee is a CLI tool that automates worktree creation and 
copies specified files to new worktrees based on configuration.

It helps manage multiple branches efficiently by creating worktrees
in .git/worktree/ directory and automatically copying configured files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sproutee - Git Worktree Management Tool")
		fmt.Println("Use 'sproutee --help' for more information.")
	},
}

var createCmd = &cobra.Command{
	Use:   "create <name> [branch]",
	Short: "Create a new worktree with file copying",
	Long: `Create a new Git worktree with the specified name and optionally 
from a specific branch. Files specified in the configuration will be 
automatically copied to the new worktree.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		branch := "HEAD"
		if len(args) > 1 {
			branch = args[1]
		}
		fmt.Printf("Creating worktree '%s' from branch '%s'\n", name, branch)
		fmt.Println("This feature is coming soon!")
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
		fmt.Println("Existing worktrees:")
		fmt.Println("This feature is coming soon!")
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up worktrees",
	Long:  "Remove unused or orphaned worktrees.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cleaning up worktrees...")
		fmt.Println("This feature is coming soon!")
	},
}

func init() {
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