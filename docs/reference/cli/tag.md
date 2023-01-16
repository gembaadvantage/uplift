# Command Line

```text
Generates a new git tag by scanning the git log of a repository for any
conventional commits since the last release (or identifiable tag). When
examining the git log, Uplift will always calculate the next semantic version
based on the most significant detected increment. Uplift automatically handles
the creation and pushing of a new git tag to the remote, but this behavior can
be disabled, to manage this action manually.

Conventional Commits is a set of rules for creating an explicit commit history,
which makes building automation tools like Uplift much easier. Uplift adheres to
v1.0.0 of the specification:

https://www.conventionalcommits.org/en/v1.0.0
```

## Usage

```text
uplift tag [flags]
```

## Examples

```text
# Tag the repository with the next calculated semantic version
uplift tag

# Identify the next semantic version and write to stdout. Repository is
# not tagged
uplift tag --next --silent

# Identify the current semantic version and write to stdout. Repository
# is not tagged
uplift tag --current

# Identify the current and next semantic versions and write both to stdout.
# Repository is not tagged
uplift tag --current --next --silent

# Ensure the calculated version explicitly aheres to the SemVer specification
# by stripping the "v" prefix from the generated tag
uplift tag --no-prefix

# Append a prerelease suffix to the next calculated semantic version
uplift tag --prerelease beta.1

# Tag the repository with the next calculated semantic version, but do not
# push the tag to the remote
uplift tag --no-push
```

## Flags

```text
    --current             output the current tag
    --fetch-all           fetch all tags from the remote repository
-h, --help                help for tag
    --next                output the next tag
    --no-prefix           strip the default 'v' prefix from the next calculated
                          semantic version
    --prerelease string   append a prerelease suffix to next calculated semantic
                          version
```

## Global Flags

```text
--config-dir string            a custom path to a directory containing uplift
                               config (default ".")
--debug                        show me everything that happens
--dry-run                      run without making any changes
--ignore-detached              ignore reported git detached HEAD error
--ignore-existing-prerelease   ignore any existing prerelease when calculating
                               next semantic version
--ignore-shallow               ignore reported git shallow clone error
--no-push                      no changes will be pushed to the git remote
--no-stage                     no changes will be git staged
--silent                       silence all logging
```
