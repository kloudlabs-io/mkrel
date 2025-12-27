package version

import (
	"testing"
)

func TestNew(t *testing.T) {
	latestTagFn := func() (string, error) { return "", nil }

	tests := []struct {
		name       string
		scheme     Scheme
		wantScheme Scheme
		wantErr    bool
	}{
		{
			name:       "calver scheme",
			scheme:     SchemeCalVer,
			wantScheme: SchemeCalVer,
		},
		{
			name:       "semver scheme",
			scheme:     SchemeSemVer,
			wantScheme: SchemeSemVer,
		},
		{
			name:    "unknown scheme",
			scheme:  "unknown",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.scheme, latestTagFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Scheme() != tt.wantScheme {
				t.Errorf("New().Scheme() = %v, want %v", got.Scheme(), tt.wantScheme)
			}
		})
	}
}

func TestParseScheme(t *testing.T) {
	tests := []struct {
		input   string
		want    Scheme
		wantErr bool
	}{
		{"calver", SchemeCalVer, false},
		{"CalVer", SchemeCalVer, false},
		{"CALVER", SchemeCalVer, false},
		{"semver", SchemeSemVer, false},
		{"SemVer", SchemeSemVer, false},
		{"SEMVER", SchemeSemVer, false},
		{"unknown", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseScheme(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseScheme(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseScheme(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
