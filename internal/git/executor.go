// Package git provides operations for interacting with git repositories.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Executor runs git commands in a specific directory.
// This is a struct - it holds data (fields) and can have methods attached.
type Executor struct {
	// Fields are like instance variables. Lowercase = private to package.
	workDir string // Directory to run commands in
	dryRun  bool   // If true, print commands instead of running them
	verbose bool   // If true, print commands before running
}

// NewExecutor creates a new Executor.
// This is a "constructor" pattern in Go - a function that returns a configured struct.
// Go doesn't have constructors, so we use New<Type> naming convention.
func NewExecutor(workDir string, dryRun, verbose bool) *Executor {
	// &Executor{...} creates a struct and returns a pointer to it.
	// Using pointers avoids copying the struct when passing it around.
	return &Executor{
		workDir: workDir,
		dryRun:  dryRun,
		verbose: verbose,
	}
}

// Run executes a git command and returns its output.
// This is a METHOD - a function attached to a type.
// The (e *Executor) part is called the "receiver" - like 'self' or 'this'.
func (e *Executor) Run(args ...string) (string, error) {
	// args ...string is a variadic parameter - accepts any number of strings
	// Called like: e.Run("status") or e.Run("checkout", "-b", "feature")

	if e.verbose || e.dryRun {
		fmt.Printf("$ git %s\n", strings.Join(args, " "))
	}

	if e.dryRun {
		return "", nil
	}

	// exec.Command creates a command to run
	cmd := exec.Command("git", args...)

	// Set the working directory
	cmd.Dir = e.workDir

	// Buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command and wait for it to complete
	err := cmd.Run()
	if err != nil {
		// Return stderr as part of the error for better debugging
		return "", fmt.Errorf("git %s failed: %w\n%s",
			strings.Join(args, " "), err, stderr.String())
	}

	// strings.TrimSpace removes leading/trailing whitespace
	return strings.TrimSpace(stdout.String()), nil
}

// RunSilent runs a command without printing, even in verbose mode.
// Useful for commands like checking if a branch exists.
func (e *Executor) RunSilent(args ...string) (string, error) {
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
