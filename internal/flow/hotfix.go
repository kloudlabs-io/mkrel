package flow

import (
	"fmt"
	"strings"

	"github.com/kloudlabs-io/mkrel/internal/version"
)

// HotfixStart begins a new hotfix.
// It creates a hotfix branch from main with a patch/hotfix version bump.
func (f *Flow) HotfixStart() error {
	f.print("==> Starting new hotfix")

	// 1. Check no hotfix already in progress
	hotfixes, err := f.repo.ListBranches("hotfix/")
	if err != nil {
		return fmt.Errorf("failed to list hotfix branches: %w", err)
	}
	if len(hotfixes) > 0 {
		return fmt.Errorf("hotfix already in progress: %s", hotfixes[0])
	}

	// 2. Get main branch
	mainBranch, err := f.repo.GetMainBranch()
	if err != nil {
		return err
	}
	f.print("    Using main branch: %s", mainBranch)

	// 3. Checkout main and ensure clean
	if err := f.repo.Checkout(mainBranch); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", mainBranch, err)
	}

	hasChanges, err := f.repo.HasUncommittedChanges()
	if err != nil {
		return err
	}
	if hasChanges {
		return fmt.Errorf("uncommitted changes in working directory")
	}

	// 4. Calculate next hotfix version
	current, err := f.versioner.Current()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	f.print("    Current version: %s", current)

	nextVersion, err := f.versioner.Next(current, version.BumpHotfix)
	if err != nil {
		return fmt.Errorf("failed to calculate next version: %w", err)
	}
	f.print("    Hotfix version: %s", nextVersion)

	// 5. Create hotfix branch
	branchName := "hotfix/" + nextVersion
	f.print("    Creating branch: %s", branchName)

	if err := f.repo.CreateBranch(branchName, mainBranch); err != nil {
		return fmt.Errorf("failed to create hotfix branch: %w", err)
	}

	f.printAlways("==> Hotfix %s started", nextVersion)
	f.printAlways("    Branch: %s", branchName)
	f.printAlways("")
	f.printAlways("    Make your fixes, then run:")
	f.printAlways("      mkrel hotfix finish")

	return nil
}

// HotfixFinish completes the current hotfix.
// It merges to main, tags, merges to develop, and pushes.
func (f *Flow) HotfixFinish() error {
	f.print("==> Finishing hotfix")

	// 1. Find hotfix branch
	hotfixes, err := f.repo.ListBranches("hotfix/")
	if err != nil {
		return fmt.Errorf("failed to list hotfix branches: %w", err)
	}
	if len(hotfixes) == 0 {
		return fmt.Errorf("no hotfix in progress")
	}
	if len(hotfixes) > 1 {
		return fmt.Errorf("multiple hotfixes in progress: %v", hotfixes)
	}

	hotfixBranch := hotfixes[0]
	f.print("    Hotfix branch: %s", hotfixBranch)

	// Extract version from branch name
	hotfixVersion := strings.TrimPrefix(hotfixBranch, "hotfix/")
	f.print("    Version: %s", hotfixVersion)

	// 2. Get main and develop branches
	mainBranch, err := f.repo.GetMainBranch()
	if err != nil {
		return err
	}
	developBranch, err := f.repo.GetDevelopBranch()
	if err != nil {
		return err
	}

	// 3. Checkout hotfix branch and verify clean
	if err := f.repo.Checkout(hotfixBranch); err != nil {
		return fmt.Errorf("failed to checkout hotfix branch: %w", err)
	}

	hasChanges, err := f.repo.HasUncommittedChanges()
	if err != nil {
		return err
	}
	if hasChanges {
		return fmt.Errorf("uncommitted changes in hotfix branch")
	}

	// 4. Merge to main
	f.print("    Merging to %s", mainBranch)
	if err := f.repo.Checkout(mainBranch); err != nil {
		return err
	}
	if err := f.repo.Merge(hotfixBranch, true); err != nil {
		return fmt.Errorf("failed to merge to %s: %w", mainBranch, err)
	}

	// 5. Create tag
	tagName, err := f.repo.FormatTag(hotfixVersion)
	if err != nil {
		return err
	}
	f.print("    Creating tag: %s", tagName)
	if err := f.repo.CreateTag(tagName, "Hotfix "+hotfixVersion); err != nil {
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

	// 8. Delete hotfix branch
	f.print("    Deleting branch: %s", hotfixBranch)
	if err := f.repo.DeleteBranch(hotfixBranch); err != nil {
		f.print("    Warning: failed to delete branch: %v", err)
	}

	f.printAlways("==> Hotfix %s released", hotfixVersion)

	return nil
}
