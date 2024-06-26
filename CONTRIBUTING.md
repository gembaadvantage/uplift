# Contributing Guidelines

The Uplift project accepts contributions via GitHub pull requests. This document outlines the process to help get your contribution accepted.

Please note we operate by a strict [code of conduct](https://github.com/gembaadvantage/uplift/blob/main/CODE_OF_CONDUCT.md) and your acceptance is required when interacting with this project.

## Getting Started

`uplift` is written using [Go 1.21+](https://go.dev/doc/install) and should be installed along with [go-task](https://taskfile.dev/#/installation), as it is preferred over using make.

Then clone `uplift`:

```sh
git clone git@github.com:gembaadvantage/uplift.git
```

`cd` into the directory and check everything is fine:

```sh
task
```

## Issues

Issues are used as the primary method for tracking anything to do with the Uplift project.

### Issue Types

There are 2 types of issue:

- `new feature`: There are used to track feature requests (or new ideas) from inception through to completion. A feature is typically raised to enhance the current functionality of the code
- `bug report`: These are used to track bugs within the code and should provide clear and concise instructions on how to replicate the issue

## Pull Requests

We use pull requests (PRs) to track any proposed code change. All PRs are merged using the `squash` strategy.

- When raising a PR it should be opened against the `main` branch of this repository
- A PR should be linked against an open issue
- The PR title should be based on the Conventional Commits standard to describe its intent
- All commits to a PR should be signed

## Conventional Commits

To ensure all commits are standardised and are supported by the build tooling, we have adopted [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). Please take a minute to familiarise yourself. It shouldn't take long.

## Signed Commits

All commits should be signed. You can sign your commits automatically by using the `git commit -S` command. Read the official Github documentation [here](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits) to set this up.
