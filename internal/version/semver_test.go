package version

import (
	"errors"
	"testing"
)

func TestSemVer_Scheme(t *testing.T) {
	sv := NewSemVer(func() (string, error) { return "", nil })
	if got := sv.Scheme(); got != SchemeSemVer {
		t.Errorf("Scheme() = %v, want %v", got, SchemeSemVer)
	}
}

func TestSemVer_Current(t *testing.T) {
	tests := []struct {
		name      string
		latestTag string
		latestErr error
		want      string
		wantErr   bool
	}{
		{
			name:      "valid tag without prefix",
			latestTag: "1.2.3",
			want:      "1.2.3",
		},
		{
			name:      "valid tag with v prefix",
			latestTag: "v1.2.3",
			want:      "1.2.3",
		},
		{
			name:      "prerelease tag",
			latestTag: "v1.2.3-rc.1",
			want:      "1.2.3-rc.1",
		},
		{
			name:      "empty tag (no releases)",
			latestTag: "",
			want:      "",
		},
		{
			name:      "error from latestTagFn",
			latestErr: errors.New("git error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv := NewSemVer(func() (string, error) {
				return tt.latestTag, tt.latestErr
			})

			got, err := sv.Current()
			if (err != nil) != tt.wantErr {
				t.Errorf("Current() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Current() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemVer_IsValid(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"1.2.3", true},
		{"0.0.1", true},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0-rc.0", true},
		{"1.0.0+build", true},
		{"1.0.0-rc.1+build", true},
		{"v1.2.3", true},          // v prefix is accepted by semver lib
		{"1.2", true},             // semver lib coerces to 1.2.0
		{"1", true},               // semver lib coerces to 1.0.0
		{"1.2.3.4", false},        // too many segments
		{"a.b.c", false},          // non-numeric
		{"", false},               // empty
		{"2025.12.25", true},      // semver lib accepts this (coerces to 2025.12.25)
	}

	sv := NewSemVer(func() (string, error) { return "", nil })

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got := sv.IsValid(tt.version); got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestSemVer_Next(t *testing.T) {
	tests := []struct {
		name    string
		current string
		bump    BumpType
		want    string
		wantErr bool
	}{
		// Minor bumps
		{
			name:    "minor bump",
			current: "1.2.3",
			bump:    BumpMinor,
			want:    "1.3.0",
		},
		{
			name:    "minor bump from zero",
			current: "0.0.0",
			bump:    BumpMinor,
			want:    "0.1.0",
		},
		// Patch bumps
		{
			name:    "patch bump",
			current: "1.2.3",
			bump:    BumpPatch,
			want:    "1.2.4",
		},
		{
			name:    "hotfix same as patch",
			current: "1.2.3",
			bump:    BumpHotfix,
			want:    "1.2.4",
		},
		// No current version
		{
			name:    "no current version minor",
			current: "",
			bump:    BumpMinor,
			want:    "0.1.0",
		},
		{
			name:    "no current version patch",
			current: "",
			bump:    BumpPatch,
			want:    "0.0.1",
		},
		// Errors
		{
			name:    "invalid current version",
			current: "invalid",
			bump:    BumpMinor,
			wantErr: true,
		},
		{
			name:    "unsupported bump type",
			current: "1.2.3",
			bump:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv := NewSemVer(func() (string, error) { return tt.current, nil })

			got, err := sv.Next(tt.current, tt.bump)
			if (err != nil) != tt.wantErr {
				t.Errorf("Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemVer_SetPrerelease(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		prerelease string
		want       string
	}{
		{
			name:       "add rc.0",
			version:    "1.2.0",
			prerelease: "rc.0",
			want:       "1.2.0-rc.0",
		},
		{
			name:       "add alpha",
			version:    "1.0.0",
			prerelease: "alpha",
			want:       "1.0.0-alpha",
		},
		{
			name:       "replace existing prerelease",
			version:    "1.2.0-beta.1",
			prerelease: "rc.0",
			want:       "1.2.0-rc.0",
		},
		{
			name:       "invalid version fallback",
			version:    "invalid",
			prerelease: "rc.0",
			want:       "invalid-rc.0",
		},
	}

	sv := NewSemVer(func() (string, error) { return "", nil })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sv.SetPrerelease(tt.version, tt.prerelease)
			if got != tt.want {
				t.Errorf("SetPrerelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemVer_RemovePrerelease(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "remove rc",
			version: "1.2.0-rc.0",
			want:    "1.2.0",
		},
		{
			name:    "remove alpha",
			version: "1.0.0-alpha.1",
			want:    "1.0.0",
		},
		{
			name:    "no prerelease",
			version: "1.2.3",
			want:    "1.2.3",
		},
		{
			name:    "invalid version with dash",
			version: "invalid-rc.0",
			want:    "invalid",
		},
		{
			name:    "invalid version without dash",
			version: "invalid",
			want:    "invalid",
		},
	}

	sv := NewSemVer(func() (string, error) { return "", nil })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sv.RemovePrerelease(tt.version)
			if got != tt.want {
				t.Errorf("RemovePrerelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemVer_IncrementPrerelease(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
		wantErr bool
	}{
		{
			name:    "increment rc.0",
			version: "1.2.0-rc.0",
			want:    "1.2.0-rc.1",
		},
		{
			name:    "increment rc.9",
			version: "1.2.0-rc.9",
			want:    "1.2.0-rc.10",
		},
		{
			name:    "increment alpha.1",
			version: "1.0.0-alpha.1",
			want:    "1.0.0-alpha.2",
		},
		{
			name:    "prerelease without number",
			version: "1.0.0-beta",
			want:    "1.0.0-beta.1",
		},
		{
			name:    "no prerelease",
			version: "1.2.3",
			wantErr: true,
		},
		{
			name:    "invalid version",
			version: "invalid",
			wantErr: true,
		},
	}

	sv := NewSemVer(func() (string, error) { return "", nil })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sv.IncrementPrerelease(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrementPrerelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("IncrementPrerelease() = %v, want %v", got, tt.want)
			}
		})
	}
}
