package version

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// SemVer implements semantic versioning (https://semver.org).
// Format: MAJOR.MINOR.PATCH with optional prerelease suffix.
type SemVer struct {
	latestTagFn func() (string, error)
}

// NewSemVer creates a SemVer versioner.
func NewSemVer(latestTagFn func() (string, error)) *SemVer {
	return &SemVer{
		latestTagFn: latestTagFn,
	}
}

// Scheme returns the versioning scheme.
func (s *SemVer) Scheme() Scheme {
	return SchemeSemVer
}

// Current returns the current version from git tags.
func (s *SemVer) Current() (string, error) {
	tag, err := s.latestTagFn()
	if err != nil {
		return "", err
	}

	// Strip "v" prefix if present
	version := strings.TrimPrefix(tag, "v")

	if version == "" {
		return "", nil
	}

	return version, nil
}

// IsValid checks if a version is valid semver.
func (s *SemVer) IsValid(version string) bool {
	_, err := semver.NewVersion(version)
	return err == nil
}

// Next calculates the next version based on bump type.
func (s *SemVer) Next(current string, bump BumpType) (string, error) {
	// If no current version, start at 0.1.0
	if current == "" {
		switch bump {
		case BumpMinor:
			return "0.1.0", nil
		case BumpPatch, BumpHotfix:
			return "0.0.1", nil
		}
	}

	// Parse current version
	v, err := semver.NewVersion(current)
	if err != nil {
		return "", fmt.Errorf("invalid current version %q: %w", current, err)
	}

	// Calculate next version
	var next semver.Version
	switch bump {
	case BumpMinor:
		// 1.2.3 -> 1.3.0
		next = v.IncMinor()
	case BumpPatch, BumpHotfix:
		// 1.2.3 -> 1.2.4
		next = v.IncPatch()
	default:
		return "", fmt.Errorf("unsupported bump type: %s", bump)
	}

	return next.String(), nil
}

// SetPrerelease adds a prerelease suffix (e.g., "1.2.0-rc.0").
func (s *SemVer) SetPrerelease(version, prerelease string) string {
	v, err := semver.NewVersion(version)
	if err != nil {
		// If invalid, just append
		return version + "-" + prerelease
	}

	// Create new version with prerelease
	newV, err := v.SetPrerelease(prerelease)
	if err != nil {
		return version + "-" + prerelease
	}

	return newV.String()
}

// RemovePrerelease removes the prerelease suffix.
func (s *SemVer) RemovePrerelease(version string) string {
	v, err := semver.NewVersion(version)
	if err != nil {
		// Try simple string manipulation
		if idx := strings.Index(version, "-"); idx != -1 {
			return version[:idx]
		}
		return version
	}

	// Create new version without prerelease
	newV, _ := v.SetPrerelease("")
	return newV.String()
}

// IncrementPrerelease increments the prerelease number.
// e.g., "1.0.0-rc.0" -> "1.0.0-rc.1"
func (s *SemVer) IncrementPrerelease(version string) (string, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return "", fmt.Errorf("invalid version: %w", err)
	}

	pre := v.Prerelease()
	if pre == "" {
		return "", fmt.Errorf("version %s has no prerelease to increment", version)
	}

	// Parse prerelease like "rc.0" or "beta.1"
	parts := strings.Split(pre, ".")
	if len(parts) < 2 {
		// Just "rc" without number, add .1
		newPre := pre + ".1"
		newV, _ := v.SetPrerelease(newPre)
		return newV.String(), nil
	}

	// Try to parse last part as number
	lastIdx := len(parts) - 1
	num := 0
	fmt.Sscanf(parts[lastIdx], "%d", &num)
	parts[lastIdx] = fmt.Sprintf("%d", num+1)

	newPre := strings.Join(parts, ".")
	newV, _ := v.SetPrerelease(newPre)
	return newV.String(), nil
}
