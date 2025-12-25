// Package cli implements the command-line interface using Cobra.
package cli

import (
	"github.com/spf13/cobra"
)

// These variables are set at build time by GoReleaser using -ldflags.
// They let us embed version info into the binary without hardcoding.
var (
	Version = "dev"     // Semantic version (e.g., "1.0.0")
	Commit  = "none"    // Git commit hash
	Date    = "unknown" // Build date
)

// rootCmd is the base command when called without any subcommands.
// In Cobra, commands are structs with fields describing their behavior.
var rootCmd = &cobra.Command{
	// Use is the one-line usage - first word becomes the command name
	Use:   "mkrel",

	// Short is a brief description shown in help listings
	Short: "Release management tool with Git Flow",

	// Long is the detailed description shown in 'mkrel --help'
	Long: `mkrel automates semantic and calendar versioning releases
following the Git Flow branching model.

It handles the complete release lifecycle:
  - Creating release/hotfix branches from develop/main
  - Bumping versions (CalVer or SemVer)
  - Merging to main and develop
  - Tagging and pushing to remote`,

	// SilenceUsage prevents printing usage on errors (cleaner output)
	SilenceUsage: true,
}

// Execute runs the root command. This is called from main().
// It returns an error if the command fails.
func Execute() error {
	return rootCmd.Execute()
}

// init() is a special Go function that runs automatically when the package loads.
// We use it to set up flags and add subcommands.
func init() {
	// Persistent flags are inherited by all subcommands
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dry-run", false, "show what would be done without making changes")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default: .mkrel.yaml)")
}
