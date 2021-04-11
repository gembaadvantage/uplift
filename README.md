# uplift

[![Build status](https://img.shields.io/github/workflow/status/gembaadvantage/uplift/ci?style=flat-square&logo=go)](https://github.com/gembaadvantage/uplift/actions?workflow=ci)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gembaadvantage/uplift?style=flat-square)](https://goreportcard.com/report/github.com/gembaadvantage/uplift)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gembaadvantage/uplift.svg?style=flat-square)](go.mod)
[![codecov](https://codecov.io/gh/gembaadvantage/uplift/branch/main/graph/badge.svg)](https://codecov.io/gh/gembaadvantage/uplift)

Semantic versioning the easy way. Automatic tagging of your respositories based on your commit messages. Powered by [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). Easy to include in your CI.

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

An initial version can be provided for the first tag. By default `0.1.0` is used. Uplift supports the use of a `v` prefix and includes it with subsequent bumps.

```sh
$ uplift bump --first v0.1.0
v0.1.0
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

0.2.0
```
