package version

import (
	"errors"
	"testing"
	"time"
)

func TestCalVer_Scheme(t *testing.T) {
	cv := NewCalVer(func() (string, error) { return "", nil })
	if got := cv.Scheme(); got != SchemeCalVer {
		t.Errorf("Scheme() = %v, want %v", got, SchemeCalVer)
	}
}

func TestCalVer_Current(t *testing.T) {
	tests := []struct {
		name        string
		latestTag   string
		latestErr   error
		want        string
		wantErr     bool
	}{
		{
			name:      "valid tag without prefix",
			latestTag: "2025.12.25",
			want:      "2025.12.25",
		},
		{
			name:      "valid tag with v prefix",
			latestTag: "v2025.12.25",
			want:      "2025.12.25",
		},
		{
			name:      "hotfix tag",
			latestTag: "2025.12.25-1",
			want:      "2025.12.25-1",
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
			cv := NewCalVer(func() (string, error) {
				return tt.latestTag, tt.latestErr
			})

			got, err := cv.Current()
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

func TestCalVer_IsValid(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"2025.12.25", true},
		{"2025.01.01", true},
		{"2025.12.25-1", true},
		{"2025.12.25-99", true},
		{"v2025.12.25", false},    // v prefix not valid
		{"2025.1.1", false},       // single digit month/day
		{"25.12.25", false},       // 2-digit year
		{"2025-12-25", false},     // wrong separator
		{"2025.12.25.1", false},   // extra segment
		{"1.2.3", false},          // semver
		{"", false},               // empty
		{"invalid", false},        // random string
	}

	cv := NewCalVer(func() (string, error) { return "", nil })

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got := cv.IsValid(tt.version); got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestCalVer_Next_Release(t *testing.T) {
	fixedTime := time.Date(2025, 12, 26, 10, 30, 0, 0, time.UTC)

	cv := &CalVer{
		latestTagFn: func() (string, error) { return "2025.12.20", nil },
		now:         func() time.Time { return fixedTime },
	}

	got, err := cv.Next("2025.12.20", BumpMinor)
	if err != nil {
		t.Fatalf("Next() error = %v", err)
	}

	want := "2025.12.26"
	if got != want {
		t.Errorf("Next() = %v, want %v", got, want)
	}
}

func TestCalVer_Next_Hotfix(t *testing.T) {
	tests := []struct {
		name    string
		current string
		today   time.Time
		bump    BumpType
		want    string
	}{
		{
			name:    "first hotfix same day",
			current: "2025.12.26",
			today:   time.Date(2025, 12, 26, 10, 0, 0, 0, time.UTC),
			bump:    BumpHotfix,
			want:    "2025.12.26-1",
		},
		{
			name:    "second hotfix same day",
			current: "2025.12.26-1",
			today:   time.Date(2025, 12, 26, 10, 0, 0, 0, time.UTC),
			bump:    BumpHotfix,
			want:    "2025.12.26-2",
		},
		{
			name:    "hotfix different day",
			current: "2025.12.25",
			today:   time.Date(2025, 12, 26, 10, 0, 0, 0, time.UTC),
			bump:    BumpHotfix,
			want:    "2025.12.26-1",
		},
		{
			name:    "hotfix from hotfix different day",
			current: "2025.12.25-3",
			today:   time.Date(2025, 12, 26, 10, 0, 0, 0, time.UTC),
			bump:    BumpHotfix,
			want:    "2025.12.26-1",
		},
		{
			name:    "BumpPatch same as BumpHotfix",
			current: "2025.12.26",
			today:   time.Date(2025, 12, 26, 10, 0, 0, 0, time.UTC),
			bump:    BumpPatch,
			want:    "2025.12.26-1",
		},
		{
			name:    "invalid current version starts fresh",
			current: "invalid",
			today:   time.Date(2025, 12, 26, 10, 0, 0, 0, time.UTC),
			bump:    BumpHotfix,
			want:    "2025.12.26-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CalVer{
				latestTagFn: func() (string, error) { return tt.current, nil },
				now:         func() time.Time { return tt.today },
			}

			got, err := cv.Next(tt.current, tt.bump)
			if err != nil {
				t.Fatalf("Next() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalVer_Next_UnsupportedBumpType(t *testing.T) {
	cv := &CalVer{
		latestTagFn: func() (string, error) { return "", nil },
		now:         func() time.Time { return time.Now() },
	}

	_, err := cv.Next("2025.12.26", "invalid")
	if err == nil {
		t.Error("Next() expected error for invalid bump type")
	}
}

func TestCalVer_SetPrerelease(t *testing.T) {
	cv := NewCalVer(func() (string, error) { return "", nil })

	// CalVer ignores prereleases
	got := cv.SetPrerelease("2025.12.26", "rc.1")
	if got != "2025.12.26" {
		t.Errorf("SetPrerelease() = %v, want %v", got, "2025.12.26")
	}
}

func TestCalVer_RemovePrerelease(t *testing.T) {
	cv := NewCalVer(func() (string, error) { return "", nil })

	// CalVer returns version unchanged
	got := cv.RemovePrerelease("2025.12.26-1")
	if got != "2025.12.26-1" {
		t.Errorf("RemovePrerelease() = %v, want %v", got, "2025.12.26-1")
	}
}

func TestCalVer_FormatForToday(t *testing.T) {
	fixedTime := time.Date(2025, 1, 5, 10, 0, 0, 0, time.UTC)

	cv := &CalVer{
		latestTagFn: func() (string, error) { return "", nil },
		now:         func() time.Time { return fixedTime },
	}

	got := cv.FormatForToday()
	want := "2025.01.05"
	if got != want {
		t.Errorf("FormatForToday() = %v, want %v", got, want)
	}
}
