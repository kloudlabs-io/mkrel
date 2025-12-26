package flow

import (
	"fmt"
	"strings"

	"github.com/kloudlabs-io/mkrel/internal/version"
)

// ReleaseStart begins a new release.
// It creates a release branch from develop with the next version.
func (f *Flow) ReleaseStart() error {
	f.print("==> Starting new release")

	// 1. Check no release already in progress
	releases, err := f.repo.ListBranches("release/")
	if err != nil {
		return fmt.Errorf("failed to list release branches: %w", err)
	}
	if len(releases) > 0 {
		return fmt.Errorf("release already in progress: %s", releases[0])
	}

	// 2. Use configured develop branch
	f.print("    Using develop branch: %s", f.devBranch)

	// 3. Checkout develop and ensure clean
	if err := f.repo.Checkout(f.devBranch); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", f.devBranch, err)
	}

	hasChanges, err := f.repo.HasUncommittedChanges()
	if err != nil {
		return err
	}
	if hasChanges {
		return fmt.Errorf("uncommitted changes in working directory")
	}

	// 4. Calculate next version
	current, err := f.versioner.Current()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	f.print("    Current version: %s", current)

	nextVersion, err := f.versioner.Next(current, version.BumpMinor)
	if err != nil {
		return fmt.Errorf("failed to calculate next version: %w", err)
	}

	// For SemVer, we might want an RC version during release
	if f.versioner.Scheme() == version.SchemeSemVer {
		nextVersion = f.versioner.SetPrerelease(nextVersion, "rc.0")
	}

	f.print("    New version: %s", nextVersion)

	// 5. Create release branch
	branchName := "release/" + nextVersion
	f.print("    Creating branch: %s", branchName)

	if err := f.repo.CreateBranch(branchName, f.devBranch); err != nil {
		return fmt.Errorf("failed to create release branch: %w", err)
	}

	f.printAlways("==> Release %s started", nextVersion)
	f.printAlways("    Branch: %s", branchName)
	f.printAlways("")
	f.printAlways("    Make any final changes, then run:")
	f.printAlways("      mkrel release finish")

	return nil
}

// ReleaseFinish completes the current release.
// It merges to main, tags, merges to develop, and pushes.
func (f *Flow) ReleaseFinish(startNew bool) error {
	f.print("==> Finishing release")

	// 1. Find release branch
	releases, err := f.repo.ListBranches("release/")
	if err != nil {
		return fmt.Errorf("failed to list release branches: %w", err)
	}
	if len(releases) == 0 {
		return fmt.Errorf("no release in progress")
	}
	if len(releases) > 1 {
		return fmt.Errorf("multiple releases in progress: %v", releases)
	}

	releaseBranch := releases[0]
	f.print("    Release branch: %s", releaseBranch)

	// Extract version from branch name (release/X.Y.Z -> X.Y.Z)
	releaseVersion := strings.TrimPrefix(releaseBranch, "release/")

	// For SemVer, remove RC suffix for final version
	finalVersion := f.versioner.RemovePrerelease(releaseVersion)
	f.print("    Final version: %s", finalVersion)

	// 2. Use configured main and develop branches
	mainBranch := f.mainBranch
	developBranch := f.devBranch

	// 3. Checkout release branch and verify clean
	if err := f.repo.Checkout(releaseBranch); err != nil {
		return fmt.Errorf("failed to checkout release branch: %w", err)
	}

	hasChanges, err := f.repo.HasUncommittedChanges()
	if err != nil {
		return err
	}
	if hasChanges {
		return fmt.Errorf("uncommitted changes in release branch")
	}

	// 4. Merge to main
	f.print("    Merging to %s", mainBranch)
	if err := f.repo.Checkout(mainBranch); err != nil {
		return err
	}
	if err := f.repo.Merge(releaseBranch, true); err != nil {
		return fmt.Errorf("failed to merge to %s: %w", mainBranch, err)
	}

	// 5. Create tag
	tagName, err := f.repo.FormatTag(finalVersion)
	if err != nil {
		return err
	}
	f.print("    Creating tag: %s", tagName)
	if err := f.repo.CreateTag(tagName, "Release "+finalVersion); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	// 6. Merge to develop
	f.print("    Merging to %s", developBranch)
	if err := f.repo.Checkout(developBranch); err != nil {
		return err
	}
	if err := f.repo.Merge(mainBranch, true); err != nil {
		return fmt.Errorf("failed to merge to %s: %w", developBranch, err)
	}

	// 7. Push everything
	f.print("    Pushing to %s", f.remote)
	if err := f.repo.PushWithTags(f.remote, mainBranch, developBranch); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	// 8. Delete release branch
	f.print("    Deleting branch: %s", releaseBranch)
	if err := f.repo.DeleteBranch(releaseBranch); err != nil {
		// Non-fatal - branch might need force delete
		f.print("    Warning: failed to delete branch: %v", err)
	}

	f.printAlways("==> Released %s", finalVersion)

	// 9. Optionally start new release
	if startNew {
		f.printAlways("")
		return f.ReleaseStart()
	}

	return nil
}
