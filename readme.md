# rtag

rtag is an interactive git tag management utility for bumping up semantic
versions.

This utility provides assistance for tagging and incrementing releases in
git-based projects in a way consistent with the semantic versioning standard (as
described in https://semver.org).

## Installation

To install a binary release:

- download the file matching your platform here: [Latest release
  binaries](https://github.com/adnsv/rtag/releases/latest)
- unzip it into the directory of your choice
- make sure your system path resolves to that directory

To build and install `rtag` from sources:

- make sure you have a recent GO compiler installed
- execute `go install github.com/adnsv/rtag@latest`

## Usage

Execute `rtag` from a command line while located in a directory checked out from
`git`.

This utility queries the state of the git repository with `git rev-parse`, `git
status`, `git describe` and other git commands.

Based on the semantic version obtained from the latest existing tag, `rtag` then
shows an interactive TUI with choices for bumping up major/minor/patch numbers,
creating alpha/beta/rc pre-releases, and more.

When a new version is selected, `rtag` creates the tag in your local repository
and offers to push it to the remote origin.

### Options

```
Usage: rtag [--prefix=<ver-prefix>]

Options:
      --version   Show the version and exit
  -p, --prefix    prefix for new tags (default "AUTO")
```

## Features

### Version Bumping

The main menu shows available version actions based on your current tag:

- **Patch** - Bump patch version (e.g., 1.0.0 → 1.0.1)
- **Minor** - Bump minor version (e.g., 1.0.0 → 1.1.0)
- **Major** - Bump major version (e.g., 1.0.0 → 2.0.0)
- **Pre-release** - Create alpha, beta, or rc releases

### Undo Last Tag

The "Undo last tag" option appears in the main menu when tags exist. Selecting
it shows:

- A list of recent tags with their local/remote status
- Preview of which tag will become the latest after deletion
- Choice of deletion scope (local only, remote only, or both)

### First Tag Creation

When no tags exist in your repository, rtag offers:

- **v0.1.0** - Recommended starting version
- **Custom version** - Enter any valid semantic version

### Dirty Repository Handling

When your repository has uncommitted changes, rtag presents options to:

- **Commit** - Stage and commit all changes before tagging
- **Stash** - Stash changes to create tag on clean state
- **Proceed anyway** - Create tag on dirty repository (links to latest commit)

## Automatic Prefixes

By default, `rtag` will match the prefix for new tags with the prefix extracted
from the last tag: `1.0.0` -> `1.1.0`, `v1.0.0` -> `v1.1.0`, `ver_1.0.0` ->
`ver_1.1.0`, etc.

You can override this behavior by specifying which prefix to use explicitly with
the `--prefix="<MYPREFIX>"` option.

## Navigation

The TUI uses standard keyboard controls:

- **↑/↓** - Navigate between options
- **Enter** - Select option
- **Esc** - Go back to previous screen
- **q** - Quit

## License

The rtag utility is licensed under the MIT license.

Other libraries used:

- https://github.com/adnsv/go-utils
- https://github.com/jawher/mow.cli
- https://github.com/blang/semver
- https://github.com/charmbracelet/bubbletea
- https://github.com/charmbracelet/bubbles
- https://github.com/charmbracelet/lipgloss
