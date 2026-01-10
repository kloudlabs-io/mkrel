package cli

import (
	"github.com/spf13/cobra"

	"github.com/kloudlabs-io/mkrel/internal/config"
	"github.com/kloudlabs-io/mkrel/internal/flow"
)

// releaseCmd is a parent command - it groups related subcommands.
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Manage release branches",
	Long: `Manage release branches following Git Flow.

A release starts from develop, allows final testing/fixes,
then merges to both main and develop with version tagging.`,
}

// releaseStartCmd starts a new release branch.
var releaseStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new release",
	Long: `Start a new release branch from develop.

This will:
  1. Verify no release is already in progress
  2. Calculate the next version (CalVer date or SemVer minor bump)
  3. Create release/<version> branch from develop`,

	RunE: runReleaseStart,
}

// releaseFinishCmd finishes the current release.
var releaseFinishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Finish the current release",
	Long: `Finish the current release branch.

This will:
  1. Finalize the version (remove RC suffix if any)
  2. Merge release branch to main
  3. Tag the release
  4. Merge back to develop
  5. Push everything to remote
  6. Delete the local release branch`,

	RunE: runReleaseFinish,
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.AddCommand(releaseStartCmd)
	releaseCmd.AddCommand(releaseFinishCmd)

}

// runReleaseStart executes the release start command.
func runReleaseStart(cmd *cobra.Command, args []string) error {
	// Get flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")
	configPath, _ := cmd.Flags().GetString("config")

	// Load config (uses defaults if no config file)
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	// Create flow with config
	f, err := flow.New(flow.Options{
		Scheme:     cfg.Scheme,
		Remote:     cfg.Remote,
		MainBranch: cfg.Branches.Main,
		DevBranch:  cfg.Branches.Develop,
		DryRun:     dryRun,
		Verbose:    verbose,
	})
	if err != nil {
		return err
	}

	return f.ReleaseStart()
}

// runReleaseFinish executes the release finish command.
func runReleaseFinish(cmd *cobra.Command, args []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")
	configPath, _ := cmd.Flags().GetString("config")

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	f, err := flow.New(flow.Options{
		Scheme:     cfg.Scheme,
		Remote:     cfg.Remote,
		MainBranch: cfg.Branches.Main,
		DevBranch:  cfg.Branches.Develop,
		DryRun:     dryRun,
		Verbose:    verbose,
	})
	if err != nil {
		return err
	}

	return f.ReleaseFinish()
}
