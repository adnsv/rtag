# rtag

rtag is an interactive git tag management utility for bumping up semantic
versions

This utility provides assistance for tagging and incrementing releases in
git-based projects in a way consistent with the semantic versioning standard (as
described in https://semver.org).

## Installation

To install a binary release:

- download the file matching your platform here: [Latest release
  binaries](https://github.com/adnsv/rtag-prev/releases/latest)
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
shows an interactive set of choices for bumping up major/minor/release numbers,
creating alpha/beta/rc pre-releases, etc. 

When a new version is selected, `rtag` then updates the tag in the local
repository; it can also push the new tag into the remote origin.

In addition, you can run the utility with the `--undo` tag to undo the latest
tagging operation in local/remote repositories.

Execute `rtag --help` to see the full list of options.

## Dirty Repositories

When running `rtag` in a repository that has uncommitted changes, the utility
shows a warning and exits. You can override this by using the `--allow-dirty`
option; in this case the tag will be linked to the latest commit.

## Automatic Prefixes

By default, `rtag` will match the prefix for new tags with the prefix extracted
from the last tag: `1.0.0` -> `1.1.0`, `v1.0.0` -> `v1.1.0`, `ver_1.0.0` ->
`ver_1.1.0`, etc.

You can override this behavior by specifying which prefix to use explicitly with
the `--prefix="<MYPREFIX>"` option.

## Example

Making public release from the release candidate:

```
$ rtag

Repository Info:
- branch:             main
- author date:        2022-10-12T12:10:21
- hash:               c2f65ce56280703dd21328d347550c91a6d563e9
- state:              clean, no uncommited changes
- auto prefix:        v
- additional commits: 1
- semantic ver:       0.6.0-rc.4+1

available actions:
1: bump 'rc' v0.6.0-rc.5
2: make release v0.6.0
type a number [1...2]:  2

Ready to tag as v0.6.0 (with comment 'tagging as v0.6.0'):

proceed [y/n]?  y

executing: 'git tag -a v0.6.0 -m "tagging as v0.6.0"

Your local repository is now tagged as v0.6.0

To push this change to remote, execute:

    git push origin v0.6.0

NOTE: to revert local and/or remote tagging,
      re-run rtag with --undo

This utility can push the new tag for you

Proceed with push [y/n]?  y

executing: 'git push origin v0.6.0'
Enumerating objects: 1, done.
Counting objects: 100% (1/1), done.
Writing objects: 100% (1/1), 164 bytes | 164.00 KiB/s, done.
Total 1 (delta 0), reused 0 (delta 0), pack-reused 0
To https://github.com/adnsv/go-utils
 * [new tag]         v0.6.0 -> v0.6.0

mission accomplished
```

Making an alpha release after that:

```
$ rtag

Repository Info:
- branch:           main
- author date:      2022-10-12T12:10:21
- hash:             c2f65ce56280703dd21328d347550c91a6d563e9
- state:            clean, no uncommited changes
- auto prefix:      v
- additional commits: 2
- semantic ver:     0.6.0+2

available actions:
1: increment patch v0.6.1 (backwards compatible bug fixes) ...
2: increment minor v0.7.0 (backwards compatible new functionality) ...
3: increment major v1.0.0 (incompatible API changes) ...
type a number [1...3]:  1

Select (pre-)release type
- 'alpha'   for v0.6.1-alpha.1
- 'beta'    for v0.6.1-beta.1
- 'rc'      for v0.6.1-rc.1
- 'release' for v0.6.1
type [alpha/beta/rc/release]: alpha

Ready to tag as v0.6.1-alpha.1 (with comment 'tagging as v0.6.1-alpha.1'):

proceed [y/n]? y

executing: 'git push origin v0.6.1-alpha.1'
Enumerating objects: 1, done.
Counting objects: 100% (1/1), done.
Writing objects: 100% (1/1), 171 bytes | 171.00 KiB/s, done.
Total 1 (delta 0), reused 0 (delta 0), pack-reused 0
To https://github.com/adnsv/go-utils
 * [new tag]         v0.6.1-alpha.1 -> v0.6.1-alpha.1

mission accomplished
```
Rolling back the previous tagging operation:

```
$ rtag --undo

Repository Info:
- branch:             main
- author date:        2022-10-12T12:10:21
- hash:               c2f65ce56280703dd21328d347550c91a6d563e9
deleting tag: v0.6.1-alpha.1

choose which tag to delete [local/remote/both]: both

ready to execute the following commands:
- 'git tag -d v0.6.1-alpha.1'
- 'git push --delete origin v0.6.1-alpha.1'

proceed [y/n]? y

deleting local tag v0.6.1-alpha.1
Deleted tag 'v0.6.1-alpha.1' (was 98171e8)
deleting remote origin tag v0.6.1-alpha.1
To https://github.com/adnsv/go-utils
 - [deleted]         v0.6.1-alpha.1

mission accomplished
```

## License

The rtag utility is licenced under the MIT license 

Other libraries used:
- https://github.com/jawher/mow.cli
- https://github.com/blang/semver
- https://github.com/josephspurrier/goversioninfo

