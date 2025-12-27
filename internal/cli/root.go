// Package cli implements the command-line interface using Cobra.
package cli

import (
	"github.com/spf13/cobra"
)

// Build-time variables set by GoReleaser via -ldflags.
var (
	Version = "dev"     // Semantic version (e.g., "1.0.0")
	Commit  = "none"    // Git commit hash
	Date    = "unknown" // Build date
)

var rootCmd = &cobra.Command{
	Use:   "mkrel",
	Short: "Release management tool with Git Flow",
	Long: `mkrel automates semantic and calendar versioning releases
following the Git Flow branching model.

It handles the complete release lifecycle:
  - Creating release/hotfix branches from develop/main
  - Bumping versions (CalVer or SemVer)
  - Merging to main and develop
  - Tagging and pushing to remote`,
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dry-run", false, "show what would be done without making changes")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default: .mkrel.yaml)")
}
