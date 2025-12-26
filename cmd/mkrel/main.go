// Package main is the entry point for the mkrel CLI application.
// In Go, the main package with a main() function is where execution begins.
package main

import (
	"os"

	"github.com/kloudlabs-io/mkrel/internal/cli"
)

func main() {
	// Execute the root command from our cli package.
	// If there's an error, exit with code 1.
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
