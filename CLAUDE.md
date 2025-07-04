# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`rtag` is an interactive Git tag management utility for semantic versioning. It helps create consistent release tags following the semver standard (https://semver.org) and provides an interactive CLI for version bumping operations.

## Development Commands

### Building
```bash
# Build for current platform
go build -o rtag

# Build with version information
go build -ldflags "-X main.app_ver=1.0.0" -o rtag

# Install from source
go install github.com/adnsv/rtag@latest
```

### Testing
**Note**: This project currently has no test suite. When adding tests:
```bash
# Run tests (when implemented)
go test ./...

# Run tests with coverage (when implemented) 
go test -cover ./...
```

### Dependencies
```bash
# Update dependencies
go mod tidy

# Download dependencies
go mod download
```

## Architecture Overview

The application follows a modular design with clear separation of concerns:

1. **main.go**: Entry point that sets up CLI using `mow.cli` and routes to either normal execution or undo mode

2. **execute.go**: Core workflow orchestrator that:
   - Retrieves repository state via `stats.go`
   - Handles version parsing and validation
   - Manages user interaction for version selection
   - Executes Git tagging operations

3. **semver.go**: Semantic versioning logic that:
   - Generates available version bump options
   - Manages pre-release progressions (alpha → beta → rc → release)
   - Handles version string manipulation

4. **stats.go**: Repository state retrieval using `adnsv/go-utils` Git operations

5. **undo.go**: Implements safe tag deletion with user confirmation for local/remote/both

6. **output.go**: Terminal-aware formatting system with ANSI escape sequence support

7. **app_version.go**: Manages application version from go.mod or build-time ldflags

## Key Workflows

### Normal Tag Creation Flow
1. Get repository statistics and validate state
2. Parse existing tags to find latest semantic version
3. Generate available version actions based on current version
4. Present interactive choices to user
5. Execute git tag command with appropriate message
6. Optionally push tag to remote origin

### Undo Flow
1. Identify latest tag
2. Ask user for deletion scope (local/remote/both)
3. Execute appropriate git commands with confirmation

## Important Implementation Details

- Uses direct Git command execution via `os/exec` rather than Git libraries
- Interactive prompts throughout for safety and user control
- Handles non-semantic tags gracefully with error recovery options
- Supports custom version prefixes (v, ver_, etc.) with auto-detection
- Validates repository cleanliness by default (override with `--allow-dirty`)
- Terminal output adapts based on capabilities (ANSI support detection)

## External Dependencies

- `github.com/adnsv/go-utils`: Git operations and utilities
- `github.com/blang/semver/v4`: Semantic version parsing
- `github.com/jawher/mow.cli`: CLI framework

## Error Handling Patterns

The codebase uses several custom error types:
- `errUserCancelled`: User declined to proceed
- `errDirtyRepo`: Repository has uncommitted changes
- `errNoSemanticTagsFound`: No semantic version tags found
- Git command failures include the full command in error messages

## User Interaction Patterns

Uses `adnsv/go-utils/prompt` for consistent interactions:
- `prompt.YN()`: Yes/No confirmations for critical operations
- `prompt.Choose()`: Numbered menu selections
- `prompt.Enum()`: String-based selections ("local"/"remote"/"both")

Confirmation required before: tag creation, push to remote, undo operations

## Git Command Execution

Commands executed via `os/exec.Command`:
- `git tag -a <tag> -m "<message>"`: Create annotated tags
- `git push origin <tag>`: Push tags to remote
- `git tag -d <tag>`: Delete local tags
- `git push --delete origin <tag>`: Delete remote tags

Commands are displayed to user before execution (dim text).

## Configuration

No configuration files or environment variables. All behavior controlled through CLI flags:
- `--prefix/-p`: Tag prefix (default: "AUTO" detection)
- `--allow-dirty/-d`: Allow tagging with uncommitted changes
- `--undo/-u`: Enter undo mode
- `--version`: Display version

## Edge Cases and Special Behaviors

- **First tag**: Defaults to `v0.1.0` (not `v0.0.1`)
- **Prefix detection**: Auto-detects from last tag (no prefix, "v", single-letter)
- **Pre-release rules**: 
  - alpha → beta/rc (not direct to release)
  - beta → rc/release
  - rc → release only
- **Version quad limit**: Warns if commits exceed 99
- **Non-semantic tags**: Falls back to older semantic tags

## Code Conventions

- Uses snake_case for functions/variables (non-idiomatic)
- Formatting functions prefixed with `fmt_`
- Print functions prefixed with `print_`
- Key-value output aligned to 18 characters
- ANSI formatting when terminal supports it

## CI/CD

GitHub Actions workflow (`.github/workflows/go.yml`) builds releases for multiple platforms:
- Linux, Windows, macOS, FreeBSD
- amd64 and arm64 architectures
- Triggered on release creation
- Version injected via ldflags during build