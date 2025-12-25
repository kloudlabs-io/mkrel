package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kloudlabs-io/mkrel/internal/config"
	"github.com/kloudlabs-io/mkrel/internal/version"
)

// initCmd creates a new .mkrel.yaml configuration file.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize mkrel configuration",
	Long: `Create a .mkrel.yaml configuration file in the current directory.

This command creates a default configuration that you can customize.
The config file controls:
  - Versioning scheme (calver or semver)
  - Branch names (main, develop)
  - Remote name
  - Optional version file updates`,

	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags for init command
	initCmd.Flags().String("scheme", "calver", "versioning scheme (calver or semver)")
	initCmd.Flags().Bool("force", false, "overwrite existing config file")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Check if config already exists
	if config.Exists() {
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			return fmt.Errorf(".mkrel.yaml already exists (use --force to overwrite)")
		}
	}

	// Parse scheme flag
	schemeStr, _ := cmd.Flags().GetString("scheme")
	scheme, err := version.ParseScheme(schemeStr)
	if err != nil {
		return err
	}

	// Create default config with specified scheme
	cfg := config.Default()
	cfg.Scheme = scheme

	// Save to file
	if err := cfg.Save(".mkrel.yaml"); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Println("Created .mkrel.yaml")
	fmt.Println("")
	fmt.Printf("  Versioning scheme: %s\n", scheme)
	fmt.Printf("  Main branch:       %s\n", cfg.Branches.Main)
	fmt.Printf("  Develop branch:    %s\n", cfg.Branches.Develop)
	fmt.Printf("  Remote:            %s\n", cfg.Remote)
	fmt.Println("")
	fmt.Println("Edit .mkrel.yaml to customize settings.")

	return nil
}
