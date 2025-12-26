package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CalVer implements calendar versioning with format YYYY.MM.DD.
// For hotfixes on the same day, it appends -1, -2, etc.
type CalVer struct {
	// latestTagFn is a function that returns the latest git tag.
	// We inject this as a function (dependency injection) rather than
	// having CalVer depend directly on the git package.
	latestTagFn func() (string, error)

	// now is the time function, injectable for testing
	now func() time.Time
}

// calverPattern matches YYYY.MM.DD or YYYY.MM.DD-N format.
// The (?:-(\d+))? part is optional, capturing the hotfix number.
var calverPattern = regexp.MustCompile(`^(\d{4})\.(\d{2})\.(\d{2})(?:-(\d+))?$`)

// NewCalVer creates a CalVer versioner.
func NewCalVer(latestTagFn func() (string, error)) *CalVer {
	return &CalVer{
		latestTagFn: latestTagFn,
		now:         time.Now, // Default to real time; can override for tests
	}
}

// Scheme returns the versioning scheme.
func (c *CalVer) Scheme() Scheme {
	return SchemeCalVer
}

// Current returns the current version from git tags.
func (c *CalVer) Current() (string, error) {
	tag, err := c.latestTagFn()
	if err != nil {
		return "", err
	}

	// Strip "v" prefix if present
	version := strings.TrimPrefix(tag, "v")

	// If no tags exist, return empty
	if version == "" {
		return "", nil
	}

	return version, nil
}

// IsValid checks if a version matches CalVer format.
func (c *CalVer) IsValid(version string) bool {
	return calverPattern.MatchString(version)
}

// Next calculates the next version.
// For releases: uses today's date (YYYY.MM.DD)
// For hotfixes: appends -N suffix (YYYY.MM.DD-1, YYYY.MM.DD-2, etc.)
func (c *CalVer) Next(current string, bump BumpType) (string, error) {
	now := c.now()
	today := fmt.Sprintf("%d.%02d.%02d", now.Year(), now.Month(), now.Day())

	switch bump {
	case BumpMinor:
		// New release: just use today's date
		return today, nil

	case BumpPatch, BumpHotfix:
		// Hotfix: need to check if we're on the same day
		return c.nextHotfix(current, today)

	default:
		return "", fmt.Errorf("unsupported bump type for CalVer: %s", bump)
	}
}

// nextHotfix calculates the next hotfix version.
func (c *CalVer) nextHotfix(current, today string) (string, error) {
	// Parse current version
	matches := calverPattern.FindStringSubmatch(current)
	if matches == nil {
		// Current version isn't valid CalVer, start fresh
		return today + "-1", nil
	}

	currentDate := fmt.Sprintf("%s.%s.%s", matches[1], matches[2], matches[3])
	hotfixNum := 0
	if matches[4] != "" {
		hotfixNum, _ = strconv.Atoi(matches[4])
	}

	if currentDate == today {
		// Same day: increment hotfix number
		return fmt.Sprintf("%s-%d", today, hotfixNum+1), nil
	}

	// Different day: new date with hotfix suffix
	return today + "-1", nil
}

// SetPrerelease is a no-op for CalVer (dates are already specific).
func (c *CalVer) SetPrerelease(version, prerelease string) string {
	// CalVer doesn't use prereleases - dates are specific enough
	return version
}

// RemovePrerelease is a no-op for CalVer.
func (c *CalVer) RemovePrerelease(version string) string {
	return version
}

// FormatForToday returns today's date as a CalVer version.
func (c *CalVer) FormatForToday() string {
	now := c.now()
	return fmt.Sprintf("%d.%02d.%02d", now.Year(), now.Month(), now.Day())
}
