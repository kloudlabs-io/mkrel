package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Repository represents a git repository and provides high-level operations.
type Repository struct {
	exec *Executor
}

// NewRepository creates a Repository for the given directory.
// If dir is empty, it uses the current working directory.
func NewRepository(dir string, dryRun, verbose bool) (*Repository, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a git repository: %s", dir)
	}

	return &Repository{
		exec: NewExecutor(dir, dryRun, verbose),
	}, nil
}

// CurrentBranch returns the name of the current branch.
func (r *Repository) CurrentBranch() (string, error) {
	return r.exec.Run("rev-parse", "--abbrev-ref", "HEAD")
}

// BranchExists checks if a branch exists (local or remote).
func (r *Repository) BranchExists(name string) bool {
	_, err := r.exec.RunSilent("show-ref", "--verify", "--quiet", "refs/heads/"+name)
	return err == nil
}

// ListBranches returns branches matching a prefix (e.g., "release/").
func (r *Repository) ListBranches(prefix string) ([]string, error) {
	output, err := r.exec.RunSilent("branch", "--list", prefix+"*")
	if err != nil {
		return nil, err
	}

	if output == "" {
		return []string{}, nil
	}

	// Parse output: each line is "  branch-name" or "* branch-name" (current)
	var branches []string
	for _, line := range strings.Split(output, "\n") {
		branch := strings.TrimSpace(line)
		branch = strings.TrimPrefix(branch, "* ") // Remove current branch marker
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

// CreateBranch creates a new branch from a base branch.
func (r *Repository) CreateBranch(name, base string) error {
	_, err := r.exec.Run("checkout", "-b", name, base)
	return err
}

// Checkout switches to the specified branch.
func (r *Repository) Checkout(branch string) error {
	_, err := r.exec.Run("checkout", branch)
	return err
}

// DeleteBranch deletes a local branch.
func (r *Repository) DeleteBranch(name string) error {
	_, err := r.exec.Run("branch", "-d", name)
	return err
}

// Merge merges a branch into the current branch.
// noFF forces a merge commit even for fast-forward merges.
func (r *Repository) Merge(branch string, noFF bool) error {
	args := []string{"merge"}
	if noFF {
		args = append(args, "--no-ff")
	}
	args = append(args, branch)

	_, err := r.exec.Run(args...)
	return err
}

// HasUncommittedChanges checks if there are uncommitted changes.
func (r *Repository) HasUncommittedChanges() (bool, error) {
	// git status --porcelain returns empty if clean
	output, err := r.exec.RunSilent("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return output != "", nil
}

// Commit creates a commit with the given message.
func (r *Repository) Commit(message string) error {
	_, err := r.exec.Run("commit", "-m", message)
	return err
}

// GetDevelopBranch finds the develop branch (might be "develop" or "development").
func (r *Repository) GetDevelopBranch() (string, error) {
	for _, name := range []string{"develop", "development", "dev"} {
		if r.BranchExists(name) {
			return name, nil
		}
	}
	return "", fmt.Errorf("no develop branch found (tried: develop, development, dev)")
}

// GetMainBranch finds the main branch (might be "main" or "master").
func (r *Repository) GetMainBranch() (string, error) {
	for _, name := range []string{"main", "master"} {
		if r.BranchExists(name) {
			return name, nil
		}
	}
	return "", fmt.Errorf("no main branch found (tried: main, master)")
}
