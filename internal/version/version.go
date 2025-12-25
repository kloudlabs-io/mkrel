// Package version handles semantic and calendar versioning.
package version

import "fmt"

// BumpType indicates what kind of version bump to perform.
type BumpType string

const (
	BumpMinor  BumpType = "minor"  // For releases (SemVer: 1.2.0 -> 1.3.0)
	BumpPatch  BumpType = "patch"  // For hotfixes (SemVer: 1.2.3 -> 1.2.4)
	BumpHotfix BumpType = "hotfix" // For CalVer hotfixes (2025.12.25 -> 2025.12.25-1)
)

// Scheme represents a versioning scheme.
type Scheme string

const (
	SchemeCalVer Scheme = "calver" // Calendar versioning (default)
	SchemeSemVer Scheme = "semver" // Semantic versioning
)

// Versioner is an INTERFACE - it defines behavior, not implementation.
// Any type that has these methods "implements" Versioner automatically.
// This is called "implicit interface implementation" - no "implements" keyword needed!
type Versioner interface {
	// Current returns the current version string.
	Current() (string, error)

	// Next calculates the next version based on bump type.
	Next(current string, bump BumpType) (string, error)

	// Scheme returns which versioning scheme this is.
	Scheme() Scheme

	// IsValid checks if a version string is valid for this scheme.
	IsValid(version string) bool

	// SetPrerelease adds a prerelease suffix (e.g., "rc.0").
	// Only applicable to SemVer; CalVer returns version unchanged.
	SetPrerelease(version, prerelease string) string

	// RemovePrerelease removes prerelease suffix.
	RemovePrerelease(version string) string
}

// New creates a Versioner for the specified scheme.
// This is a FACTORY FUNCTION - returns an interface, not a concrete type.
// Callers don't need to know if they're getting CalVer or SemVer.
func New(scheme Scheme, latestTagFn func() (string, error)) (Versioner, error) {
	switch scheme {
	case SchemeCalVer:
		return NewCalVer(latestTagFn), nil
	case SchemeSemVer:
		return NewSemVer(latestTagFn), nil
	default:
		return nil, fmt.Errorf("unknown versioning scheme: %s", scheme)
	}
}

// ParseScheme converts a string to a Scheme.
func ParseScheme(s string) (Scheme, error) {
	switch s {
	case "calver", "CalVer", "CALVER":
		return SchemeCalVer, nil
	case "semver", "SemVer", "SEMVER":
		return SchemeSemVer, nil
	default:
		return "", fmt.Errorf("unknown scheme: %s (use 'calver' or 'semver')", s)
	}
}
