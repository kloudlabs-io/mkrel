package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kloudlabs-io/mkrel/internal/version"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Scheme != version.SchemeCalVer {
		t.Errorf("Default().Scheme = %v, want %v", cfg.Scheme, version.SchemeCalVer)
	}
	if cfg.CalVerFormat != "YYYY.MM.DD" {
		t.Errorf("Default().CalVerFormat = %v, want %v", cfg.CalVerFormat, "YYYY.MM.DD")
	}
	if cfg.Branches.Main != "main" {
		t.Errorf("Default().Branches.Main = %v, want %v", cfg.Branches.Main, "main")
	}
	if cfg.Branches.Develop != "develop" {
		t.Errorf("Default().Branches.Develop = %v, want %v", cfg.Branches.Develop, "develop")
	}
	if cfg.Remote != "origin" {
		t.Errorf("Default().Remote = %v, want %v", cfg.Remote, "origin")
	}
}

func TestLoad_NoConfigFile(t *testing.T) {
	// Create a temp directory with no config file
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Should return defaults
	if cfg.Scheme != version.SchemeCalVer {
		t.Errorf("Load().Scheme = %v, want %v", cfg.Scheme, version.SchemeCalVer)
	}
}

func TestLoad_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".mkrel.yaml")

	configContent := `
scheme: semver
calver_format: "YY.MM.DD"
branches:
  main: production
  develop: development
remote: upstream
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Scheme != version.SchemeSemVer {
		t.Errorf("Load().Scheme = %v, want %v", cfg.Scheme, version.SchemeSemVer)
	}
	if cfg.Branches.Main != "production" {
		t.Errorf("Load().Branches.Main = %v, want %v", cfg.Branches.Main, "production")
	}
	if cfg.Branches.Develop != "development" {
		t.Errorf("Load().Branches.Develop = %v, want %v", cfg.Branches.Develop, "development")
	}
	if cfg.Remote != "upstream" {
		t.Errorf("Load().Remote = %v, want %v", cfg.Remote, "upstream")
	}
}

func TestLoad_PartialConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".mkrel.yaml")

	// Only override some values
	configContent := `
scheme: semver
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Overridden value
	if cfg.Scheme != version.SchemeSemVer {
		t.Errorf("Load().Scheme = %v, want %v", cfg.Scheme, version.SchemeSemVer)
	}

	// Default values should be preserved
	if cfg.Branches.Main != "main" {
		t.Errorf("Load().Branches.Main = %v, want %v", cfg.Branches.Main, "main")
	}
	if cfg.Remote != "origin" {
		t.Errorf("Load().Remote = %v, want %v", cfg.Remote, "origin")
	}
}

func TestLoad_InvalidScheme(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".mkrel.yaml")

	configContent := `
scheme: invalid
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Load() expected error for invalid scheme")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".mkrel.yaml")

	configContent := `
scheme: [invalid yaml
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Load() expected error for invalid YAML")
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// No config file exists
	if Exists() {
		t.Error("Exists() = true, want false")
	}

	// Create config file
	if err := os.WriteFile(".mkrel.yaml", []byte("scheme: calver"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Now it should exist
	if !Exists() {
		t.Error("Exists() = false, want true")
	}
}

func TestFindConfigFile(t *testing.T) {
	// Create a temp directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub", "dir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirs: %v", err)
	}

	// Create config in root
	configPath := filepath.Join(tmpDir, ".mkrel.yaml")
	if err := os.WriteFile(configPath, []byte("scheme: calver"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Change to subdirectory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(subDir)

	// Should find config in parent directory
	found, err := FindConfigFile()
	if err != nil {
		t.Fatalf("FindConfigFile() error = %v", err)
	}

	if found != configPath {
		t.Errorf("FindConfigFile() = %v, want %v", found, configPath)
	}
}

func TestFindConfigFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// No config file anywhere
	found, err := FindConfigFile()
	if err != nil {
		t.Fatalf("FindConfigFile() error = %v", err)
	}

	if found != "" {
		t.Errorf("FindConfigFile() = %v, want empty string", found)
	}
}

func TestConfig_Save(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".mkrel.yaml")

	cfg := &Config{
		Scheme:       version.SchemeSemVer,
		CalVerFormat: "YYYY.MM.DD",
		Branches: BranchConfig{
			Main:    "production",
			Develop: "development",
		},
		Remote: "upstream",
	}

	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Save() did not create file")
	}

	// Load it back
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Scheme != cfg.Scheme {
		t.Errorf("Loaded.Scheme = %v, want %v", loaded.Scheme, cfg.Scheme)
	}
	if loaded.Branches.Main != cfg.Branches.Main {
		t.Errorf("Loaded.Branches.Main = %v, want %v", loaded.Branches.Main, cfg.Branches.Main)
	}
}
