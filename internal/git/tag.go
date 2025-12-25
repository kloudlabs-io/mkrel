package git

import (
	"fmt"
	"sort"
	"strings"
)

// CreateTag creates an annotated tag with a message.
func (r *Repository) CreateTag(name, message string) error {
	_, err := r.exec.Run("tag", "-a", name, "-m", message)
	return err
}

// TagExists checks if a tag exists.
func (r *Repository) TagExists(name string) bool {
	_, err := r.exec.RunSilent("show-ref", "--verify", "--quiet", "refs/tags/"+name)
	return err == nil
}

// LatestTag returns the most recent tag.
// Returns empty string if no tags exist.
func (r *Repository) LatestTag() (string, error) {
	// git describe --tags --abbrev=0 gets the most recent tag
	output, err := r.exec.RunSilent("describe", "--tags", "--abbrev=0")
	if err != nil {
		// No tags exist - this is not an error for our use case
		if strings.Contains(err.Error(), "No names found") ||
			strings.Contains(err.Error(), "No tags") {
			return "", nil
		}
		return "", err
	}
	return output, nil
}

// ListTags returns all tags, optionally filtered by prefix.
func (r *Repository) ListTags(prefix string) ([]string, error) {
	args := []string{"tag", "--list"}
	if prefix != "" {
		args = append(args, prefix+"*")
	}

	output, err := r.exec.RunSilent(args...)
	if err != nil {
		return nil, err
	}

	if output == "" {
		return []string{}, nil
	}

	tags := strings.Split(output, "\n")
	// Sort tags (git doesn't guarantee order)
	sort.Strings(tags)
	return tags, nil
}

// Push pushes refs (branches, tags) to a remote.
func (r *Repository) Push(remote string, refs ...string) error {
	args := append([]string{"push", remote}, refs...)
	_, err := r.exec.Run(args...)
	return err
}

// PushWithTags pushes refs and all tags to a remote.
func (r *Repository) PushWithTags(remote string, refs ...string) error {
	args := append([]string{"push", "--follow-tags", remote}, refs...)
	_, err := r.exec.Run(args...)
	return err
}

// FetchTags fetches all tags from a remote.
func (r *Repository) FetchTags(remote string) error {
	_, err := r.exec.Run("fetch", "--tags", remote)
	return err
}

// GetTagsOnCommit returns tags pointing to a specific commit.
func (r *Repository) GetTagsOnCommit(commit string) ([]string, error) {
	output, err := r.exec.RunSilent("tag", "--points-at", commit)
	if err != nil {
		return nil, err
	}

	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}

// GetCurrentTags returns tags pointing to HEAD.
func (r *Repository) GetCurrentTags() ([]string, error) {
	return r.GetTagsOnCommit("HEAD")
}

// DeleteTag deletes a local tag.
func (r *Repository) DeleteTag(name string) error {
	_, err := r.exec.Run("tag", "-d", name)
	return err
}

// VersionTagPrefix returns the prefix used for version tags (e.g., "v").
// This checks existing tags to determine the pattern.
func (r *Repository) VersionTagPrefix() (string, error) {
	tags, err := r.ListTags("")
	if err != nil {
		return "", err
	}

	// Check if tags use "v" prefix
	vCount := 0
	noVCount := 0
	for _, tag := range tags {
		if strings.HasPrefix(tag, "v") && len(tag) > 1 {
			// Check if it looks like a version (v followed by digit)
			if tag[1] >= '0' && tag[1] <= '9' {
				vCount++
			}
		} else if len(tag) > 0 && tag[0] >= '0' && tag[0] <= '9' {
			noVCount++
		}
	}

	// Use "v" prefix if that's the dominant pattern, or if no tags exist
	if vCount >= noVCount {
		return "v", nil
	}
	return "", nil
}

// FormatTag formats a version string with the appropriate prefix.
func (r *Repository) FormatTag(version string) (string, error) {
	prefix, err := r.VersionTagPrefix()
	if err != nil {
		return "", fmt.Errorf("failed to determine tag prefix: %w", err)
	}

	// Don't double-prefix
	if strings.HasPrefix(version, "v") && prefix == "v" {
		return version, nil
	}

	return prefix + version, nil
}
