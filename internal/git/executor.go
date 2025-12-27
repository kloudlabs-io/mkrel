// Package git provides operations for interacting with git repositories.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Executor runs git commands in a specific directory.
type Executor struct {
	workDir string
	dryRun  bool
	verbose bool
}

// NewExecutor creates a new Executor.
func NewExecutor(workDir string, dryRun, verbose bool) *Executor {
	return &Executor{
		workDir: workDir,
		dryRun:  dryRun,
		verbose: verbose,
	}
}

// Run executes a git command and returns its output.
func (e *Executor) Run(args ...string) (string, error) {
	if e.verbose || e.dryRun {
		fmt.Printf("$ git %s\n", strings.Join(args, " "))
	}

	if e.dryRun {
		return "", nil
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = e.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w\n%s",
			strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// RunSilent runs a command without printing, even in verbose mode.
// Useful for read-only commands like checking if a branch exists.
// Note: This always executes, even in dry-run mode, because it's used
// for read-only queries that don't modify the repository.
func (e *Executor) RunSilent(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = e.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w\n%s",
			strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// RunWithInput runs a git command with stdin input.
// Used for commands that need input, like commit with message from stdin.
func (e *Executor) RunWithInput(input string, args ...string) (string, error) {
	if e.verbose || e.dryRun {
		fmt.Printf("$ git %s\n", strings.Join(args, " "))
	}

	if e.dryRun {
		return "", nil
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = e.workDir
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w\n%s",
			strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
