package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command.
// It prints the version, commit, and build date.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	// Run is the function that executes when the command is called.
	// It receives the command itself and any positional arguments.
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mkrel %s\n", Version)
		fmt.Printf("  commit: %s\n", Commit)
		fmt.Printf("  built:  %s\n", Date)
	},
}

// init adds this command to the root command.
// Each file can have its own init() - they all run on package load.
func init() {
	rootCmd.AddCommand(versionCmd)
}
