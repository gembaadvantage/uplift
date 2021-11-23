# uplift

[![Build status](https://img.shields.io/github/workflow/status/gembaadvantage/uplift/ci?style=flat-square&logo=go)](https://github.com/gembaadvantage/uplift/actions?workflow=ci)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gembaadvantage/uplift?style=flat-square)](https://goreportcard.com/report/github.com/gembaadvantage/uplift)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gembaadvantage/uplift.svg?style=flat-square)](go.mod)
[![codecov](https://codecov.io/gh/gembaadvantage/uplift/branch/main/graph/badge.svg)](https://codecov.io/gh/gembaadvantage/uplift)

Semantic versioning the easy way. Automatic tagging and version bumping of files in your repositories based on your commit messages. Powered by [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). Easy to include in your CI.

:octocat: Github [action](https://github.com/marketplace/actions/uplift-action) available.

## Install

Binary downloads of `uplift` can be found on the [Releases](https://github.com/gembaadvantage/uplift/releases) page. Unpack the `uplift` binary and add it to your PATH.

### Homebrew

To use [Homebrew](https://brew.sh/):

```sh
brew tap gembaadvantage/tap
brew install uplift
```

### GoFish

To use [Fish](https://gofi.sh/):

```sh
gofish install uplift
```

### Scoop

To use [Scoop](https://scoop.sh/):

```sh
scoop install uplift
```

### Script

To install using a shell script:

```sh
curl https://raw.githubusercontent.com/gembaadvantage/uplift/master/scripts/install > install
chmod 700 install
./install
```

## Quick Start

Uplift can carry out different semantic versioning operations on your repository. All operations support the following global flags:

- `--dry-run` provide a preview of changes only. Nothing is changed
- `--debug` log extra output to the console in debug mode
- `--no-push` no changes will be committed and pushed back to your git remote
- `--silent` silence all log output from uplift

Uplift supports the use of a `v` prefix and includes it with subsequent bumps.

### Tagging

Based on the latest commit, the repository will be tagged with the next semantic version.

```sh
$ uplift tag

  • latest commit
      • retrieved latest commit   author=joe.bloggs email=joe.bloggs@example.com message=feat: new feature
   • current version
      • identified version        current=0.1.2
   • next version
      • identified next version   commit=feat: new feature current=0.1.2 next=0.2.0
   • next commit
      • changes will be committed with email=joe.bloggs@example.com message=ci(uplift): uplifted for version 0.2.0 name=joe.bloggs
   • git tag
      • identified next tag       tag=0.2.0
      • tagged repository with standard tag
```

#### Next Tag Only

Identify the next tag without tagging a repository by using the `--next` flag:

```sh
$ uplift tag --next --silent
0.2.0
```

### File Bumping

When configured, the version within any file in a git repository can be bumped to the next semantic version. The version is identified by inspecting the latest commit.

```sh
$ uplift bump

  • latest commit
      • retrieved latest commit   author=joe.bloggs email=joe.bloggs@example.com message=feat: new feature
   • current version
      • identified version        current=1.0.0
   • next version
      • identified next version   commit=feat: new feature current=1.0.0 next=1.1.0
   • next commit
      • changes will be committed with email=joe.bloggs@example.com message=ci(uplift): uplifted for version 1.1.0 name=joe.bloggs
   • bump
      • file bumped               current=1.0.0 file=chart/test/Chart.yaml next=1.1.0
      • successfully staged file  file=chart/test/Chart.yaml
   • git push
      • commit outstanding changes
      • push commit to remote
```

### Changelog

A changelog can be generated for the latest tagged semantic release and written to a `CHANGELOG.md` file. File will be created if one doesn't exist. Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

```sh
$ uplift changelog

   • latest commit
      • retrieved latest commit   author=joe.bloggs email=joe.bloggs@example.com message=fix: a bug fix
   • next commit
      • changes will be committed with email=paul.t@gembaadvantage.com message=ci(uplift): uplifted for version 1.0.1 name=paul.t
   • changelog
      • determine changes for release tag=1.0.1
      • changeset identified      commits=3 date=2021-11-19 tag=1.0.1
   • git push
      • commit outstanding changes
      • push commit to remote
```

#### Diff Only

Output the changelog diff without modifying an existing `CHANGELOG.md` by using the `--diff-only` flag:

```sh
$ uplift changelog --diff-only --silent
## [1.2.0] - 2021-11-23

1c85055 feat: a brand new feature
```

### Release

A full semantic release will be carried out. Combining both the `bump` and `tag` operations, in that order.

```sh
uplift release
```

#### Check Release

Check if uplift will carry out a release without running the release workflow:

```sh
uplift release --check
```

## Configuration

Uplift can be configured through the existence of an optional `.uplift.yml` configuration file in the root of your repository. Other supported variations are: `.uplift.yaml`, `uplift.yml` and `uplift.yaml`.

```yaml
# An initial version that will be used for the first tag in your repository.
# Tags with a 'v' prefix are supported. You cannot change format after the first tag
# Defaults to 0.1.0
firstVersion: v1.0.0

# Use annotated tags instead of lightweight tags when tagging a version bump. An
# annotated tag is treated like a regular commit and contains both author details
# and a commit message. Uses the same commit message and author details provided
# Defaults to false
annotatedTags: true

# Changes the  default commit message used when committing any staged changes
# Defaults to ci(uplift): uplifted for version $VERSION
commitMessage: "chore: a custom commit message"

# Changes the commit author used when committing any staged changes
commitAuthor:
  # Name of the author
  # Defaults to the author name within the last commit
  name: "joe.bloggs"

  # Email of the author
  # Defaults to the author email within the last commit
  email: "joe.bloggs@gmail.com"

# A list of files whose version numbers should be bumped and kept in sync with the
# latest calculated repository tag.
# Defaults to an empty list
bumps:
  - # The path of the file relative to where uplift is executed
    file: ./chart/my-chart/Chart.yaml

    # A regex for matching a version within the file
    regex: "version: $VERSION"

    # If the matched version in the file should be replaced with a semantic version.
    # This will strip any 'v' prefix if needed
    # Defaults to false
    semver: true

    # The number of times any matched version should be replaced
    # Defaults to 0, which replaces all matches
    count: 1
```

### $VERSION

**`$VERSION`** is a placeholder and will match any semantic version, including a version with an optional `v` prefix.

## Transient Tags

Uplift automatically ignores any transient tag when identifying the latest version. This allows multiple tagging conventions to be safely used on your repository. A transient tag, is any tag that is not in the `1.0.0` or `v1.0.0` version format.
