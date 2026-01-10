# mkrel

Release management tool with Git Flow and CalVer/SemVer support.

## Features

- **Calendar Versioning (CalVer)** - Default versioning using `YYYY.MM.DD` format
- **Semantic Versioning (SemVer)** - Optional traditional `MAJOR.MINOR.PATCH` format
- **Git Flow integration** - Automated release and hotfix branch workflows
- **No dependencies** - Single binary, no git-flow CLI required
- **Cross-platform** - macOS, Linux, and Windows support

## Installation

### Homebrew (macOS/Linux)

```shell
brew tap kloudlabs-io/tap
brew install mkrel
```

### Binary Download

Download the latest release from [GitHub Releases](https://github.com/kloudlabs-io/mkrel/releases).

### From Source

```shell
go install github.com/kloudlabs-io/mkrel/cmd/mkrel@latest
```

## Quick Start

```shell
# Initialize configuration (optional)
mkrel init

# Start a release from develop
mkrel release start

# Make changes, then finish the release
mkrel release finish

# Start a hotfix from main
mkrel hotfix start

# Finish the hotfix
mkrel hotfix finish
```

## Commands

### mkrel release start

Creates a new release branch from develop with the next version:

- CalVer: Uses today's date (e.g., `2025.12.25`)
- SemVer: Bumps minor version (e.g., `1.2.0` → `1.3.0-rc.0`)

### mkrel release finish

Finishes the current release:

1. Merges release branch to main
2. Tags the release
3. Merges back to develop
4. Pushes everything to remote

### mkrel hotfix start

Creates a hotfix branch from main with a patch version:

- CalVer: Appends suffix (e.g., `2025.12.25-1`)
- SemVer: Bumps patch version (e.g., `1.2.3` → `1.2.4`)

### mkrel hotfix finish

Finishes the hotfix (same flow as release finish).

### mkrel init

Creates a `.mkrel.yaml` configuration file with defaults.

## Configuration

Create `.mkrel.yaml` in your repository root:

```yaml
# Versioning scheme: calver (default) or semver
scheme: calver

# CalVer format
calver_format: YYYY.MM.DD

# Branch names
branches:
  main: main
  develop: develop

# Git remote
remote: origin
```

## Global Flags

- `--dry-run` - Show what would happen without making changes
- `-v, --verbose` - Verbose output
- `-c, --config` - Path to config file

## License

MIT License - Copyright (c) 2020-2024 Sergei Kolobov, 2025 KloudLabs LLC
