// Package flow implements Git Flow release and hotfix workflows.
package flow

import (
	"fmt"

	"github.com/kloudlabs-io/mkrel/internal/git"
	"github.com/kloudlabs-io/mkrel/internal/version"
)

// Flow orchestrates Git Flow operations for releases and hotfixes.
type Flow struct {
	repo       *git.Repository
	versioner  version.Versioner
	remote     string // Remote name (usually "origin")
	mainBranch string // Main/production branch name
	devBranch  string // Development branch name
	dryRun     bool
	verbose    bool
}

// Options configures a Flow instance.
type Options struct {
	WorkDir    string         // Repository directory (empty = current)
	Scheme     version.Scheme // Versioning scheme
	Remote     string         // Git remote name
	MainBranch string         // Main/production branch name (empty = auto-detect)
	DevBranch  string         // Development branch name (empty = auto-detect)
	DryRun     bool
	Verbose    bool
}

// New creates a new Flow instance.
func New(opts Options) (*Flow, error) {
	// Create repository wrapper
	repo, err := git.NewRepository(opts.WorkDir, opts.DryRun, opts.Verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Create versioner with a function to get latest tag
	// This is dependency injection: versioner doesn't depend on git package
	latestTagFn := func() (string, error) {
		return repo.LatestTag()
	}

	versioner, err := version.New(opts.Scheme, latestTagFn)
	if err != nil {
		return nil, err
	}

	remote := opts.Remote
	if remote == "" {
		remote = "origin"
	}

	// Use configured branches or auto-detect
	mainBranch := opts.MainBranch
	if mainBranch == "" {
		mainBranch, err = repo.GetMainBranch()
		if err != nil {
			return nil, err
		}
	}

	devBranch := opts.DevBranch
	if devBranch == "" {
		devBranch, err = repo.GetDevelopBranch()
		if err != nil {
			return nil, err
		}
	}

	return &Flow{
		repo:       repo,
		versioner:  versioner,
		remote:     remote,
		mainBranch: mainBranch,
		devBranch:  devBranch,
		dryRun:     opts.DryRun,
		verbose:    opts.Verbose,
	}, nil
}

// print outputs a message, respecting verbose mode.
func (f *Flow) print(format string, args ...interface{}) {
	// Always print in dry-run, otherwise respect verbose
	if f.dryRun || f.verbose {
		fmt.Printf(format+"\n", args...)
	}
}

// printAlways outputs a message regardless of verbose mode.
func (f *Flow) printAlways(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
