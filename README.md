# uplift

[![Build status](https://img.shields.io/github/workflow/status/gembaadvantage/uplift/ci?style=flat-square&logo=go)](https://github.com/gembaadvantage/uplift/actions?workflow=ci)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gembaadvantage/uplift?style=flat-square)](https://goreportcard.com/report/github.com/gembaadvantage/uplift)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gembaadvantage/uplift.svg?style=flat-square)](go.mod)
[![codecov](https://codecov.io/gh/gembaadvantage/uplift/branch/main/graph/badge.svg)](https://codecov.io/gh/gembaadvantage/uplift)

Semantic versioning the easy way. Automatic tagging and version bumping of files in your respositories based on your commit messages. Powered by [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). Easy to include in your CI.

## Install

Binary downloads of Uplft can be found on the [Releases](https://github.com/gembaadvantage/uplift/releases) page. Unpack the `uplift` binary and add it to your PATH.

### Homebrew

To use [Homebrew](https://brew.sh/):

```sh
brew tap gembaadvantage/tap
brew install uplift
```

## Quick Start

Based on the latest commit, the repository will be tagged with the next calculated version.

```sh
$ uplift bump
0.2.0
```

Uplift supports the use of a `v` prefix and includes it with subsequent bumps.

```sh
$ uplift bump
v0.2.0
```

A `dry run` can be carried out with optional `verbose` output, to show what uplift is up to.

```sh
$ uplift bump --dry-run --verbose

✅ git repo found
✅ retrieved latest commit:
'feat: a new snazzy feature'
✅ commit contains a bump prefix, increment identified as 'Minor'
ℹ️ existing version found: 0.1.0
✅ bumped version to: 0.2.0
ℹ️ no files to bump, skipping!
0.2.0
```

## Configuration

Uplift can be configured through the existance of an optional `.uplift.yml` configuration file in the root of your repository. Other supported variations are: `.uplift.yaml`, `uplift.yml` and `uplift.yaml`.

```yaml
firstVersion: v1.0.0
bumps:
  - file: ./chart/my-chart/Chart.yaml
    regex: "version: $VERSION"
    semver: true
    count: 1
```

| Option       | Required | Default         | Description                                                                                            |
| ------------ | -------- | --------------- | ------------------------------------------------------------------------------------------------------ |
| firstVersion | No       | 0.1.0           | An initial version for the first tag against the repository                                            |
| bumps        | No       | []              | A list of files whose version numbers should be bumped and kept in sync with the latest repository tag |
| bumps.file   | Yes      | -               | The path of the file, relative to the root of the repository                                           |
| bumps.regex  | Yes      | -               | A regex used to identify a version within the file that should be bumped                               |
| bumps.semver | No       | false           | The version in the file must be a semantic version, ensuring any `v` prefix is removed                 |
| bumps.count  | No       | 0 (All Matches) | The number of times the matching version should be bumped within the file                              |

### $VERSION

**`$VERSION`** is a placeholder and will match any semantic version, including a version with an optional `v` prefix.

## Transient Tags

Uplift automatically ignores any transient tag when identifying the latest version. This allows multiple tagging conventions to be safely used on your repository. A transient tag, is any tag that is not in the `1.0.0` or `v1.0.0` version format.
