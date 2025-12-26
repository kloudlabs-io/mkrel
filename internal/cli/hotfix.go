package cli

import (
	"github.com/spf13/cobra"

	"github.com/kloudlabs-io/mkrel/internal/config"
	"github.com/kloudlabs-io/mkrel/internal/flow"
)

// hotfixCmd groups hotfix-related subcommands.
var hotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "Manage hotfix branches",
	Long: `Manage hotfix branches following Git Flow.

A hotfix starts from main (production), applies urgent fixes,
then merges to both main and develop with version tagging.`,
}

// hotfixStartCmd starts a new hotfix branch.
var hotfixStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new hotfix",
	Long: `Start a new hotfix branch from main.

This will:
  1. Verify no hotfix is already in progress
  2. Calculate the next hotfix version
  3. Create hotfix/<version> branch from main`,

	RunE: runHotfixStart,
}

// hotfixFinishCmd finishes the current hotfix.
var hotfixFinishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Finish the current hotfix",
	Long: `Finish the current hotfix branch.

This will:
  1. Merge hotfix branch to main
  2. Tag the hotfix release
  3. Merge back to develop
  4. Push everything to remote
  5. Delete the local hotfix branch`,

	RunE: runHotfixFinish,
}

func init() {
	rootCmd.AddCommand(hotfixCmd)
	hotfixCmd.AddCommand(hotfixStartCmd)
	hotfixCmd.AddCommand(hotfixFinishCmd)
}

// runHotfixStart executes the hotfix start command.
func runHotfixStart(cmd *cobra.Command, args []string) error {
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

	return f.HotfixStart()
}

// runHotfixFinish executes the hotfix finish command.
func runHotfixFinish(cmd *cobra.Command, args []string) error {
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

	return f.HotfixFinish()
}
