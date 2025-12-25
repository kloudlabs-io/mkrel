// Package config handles loading and managing mkrel configuration.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/kloudlabs-io/mkrel/internal/version"
)

// Config holds all configuration for mkrel.
type Config struct {
	// Scheme is the versioning scheme: "calver" or "semver"
	Scheme version.Scheme `mapstructure:"scheme"`

	// CalVerFormat is the CalVer format (default: "YYYY.MM.DD")
	CalVerFormat string `mapstructure:"calver_format"`

	// Branches configures branch names
	Branches BranchConfig `mapstructure:"branches"`

	// Remote is the git remote name (default: "origin")
	Remote string `mapstructure:"remote"`

	// VersionFiles lists files to update with version (optional)
	VersionFiles []VersionFile `mapstructure:"version_files"`
}

// BranchConfig holds branch naming configuration.
type BranchConfig struct {
	Main    string `mapstructure:"main"`    // Production branch (default: "main")
	Develop string `mapstructure:"develop"` // Development branch (default: "develop")
}

// VersionFile describes a file to update with version info.
type VersionFile struct {
	Path    string `mapstructure:"path"`    // File path
	Pattern string `mapstructure:"pattern"` // Pattern with {{version}} placeholder
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		Scheme:       version.SchemeCalVer,
		CalVerFormat: "YYYY.MM.DD",
		Branches: BranchConfig{
			Main:    "main",
			Develop: "develop",
		},
		Remote:       "origin",
		VersionFiles: []VersionFile{},
	}
}

// Load reads configuration from file and environment.
// It looks for .mkrel.yaml in the current directory.
func Load(configPath string) (*Config, error) {
	// Start with defaults
	cfg := Default()

	// Set up Viper
	v := viper.New()

	// Set config file name and type
	if configPath != "" {
		// Explicit config file path
		v.SetConfigFile(configPath)
	} else {
		// Look for .mkrel.yaml in current directory
		v.SetConfigName(".mkrel")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	}

	// Set defaults in Viper (these become the fallbacks)
	v.SetDefault("scheme", string(cfg.Scheme))
	v.SetDefault("calver_format", cfg.CalVerFormat)
	v.SetDefault("branches.main", cfg.Branches.Main)
	v.SetDefault("branches.develop", cfg.Branches.Develop)
	v.SetDefault("remote", cfg.Remote)

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		// Config file not found is OK - use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Other errors are real problems
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal into our struct
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Parse scheme string into type
	if schemeStr := v.GetString("scheme"); schemeStr != "" {
		scheme, err := version.ParseScheme(schemeStr)
		if err != nil {
			return nil, err
		}
		cfg.Scheme = scheme
	}

	return cfg, nil
}

// Save writes the configuration to a file.
func (c *Config) Save(path string) error {
	v := viper.New()

	v.Set("scheme", string(c.Scheme))
	v.Set("calver_format", c.CalVerFormat)
	v.Set("branches.main", c.Branches.Main)
	v.Set("branches.develop", c.Branches.Develop)
	v.Set("remote", c.Remote)

	if len(c.VersionFiles) > 0 {
		v.Set("version_files", c.VersionFiles)
	}

	return v.WriteConfigAs(path)
}

// Exists checks if a config file exists in the current directory.
func Exists() bool {
	_, err := os.Stat(".mkrel.yaml")
	return err == nil
}

// FindConfigFile looks for config file in current directory and parents.
func FindConfigFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		configPath := filepath.Join(dir, ".mkrel.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, no config found
			return "", nil
		}
		dir = parent
	}
}
