# Command Line

```text
Scans the git log for the latest semantic release and generates a changelog
entry. If this is a first release, all commits between the last release (or
identifiable tag) and the repository trunk will be written to the changelog.
Any subsequent entry within the changelog will only contain commits between
the latest set of tags. Basic customization is supported. Optionally commits
can be explicitly included or excluded from the entry and sorted in ascending
or descending order. Uplift automatically handles the staging and pushing of
changes to the CHANGELOG.md file to the git remote, but this behavior can be
disabled, to manage this action manually.

Uplift bases its changelog format on the Keep a Changelog specification:

https://keepachangelog.com/en/1.0.0/
```

## Usage

```text
uplift changelog [flags]
```

## Examples

```text
# Generate the next changelog entry for the latest semantic release
uplift changelog

# Generate a changelog for the entire history of the repository
uplift changelog --all

# Generate the next changelog entry and write it to stdout
uplift changelog --diff-only

# Generate the next changelog entry by exclude any conventional commits
# with the ci, chore or test prefixes
uplift changelog --exclude "^ci,^chore,^test"

# Generate the next changelog entry with commits that only include the
# following scope
uplift changelog --include "^.*\(scope\)"

# Generate the next changelog entry but do not stage or push any changes
# back to the git remote
uplift changelog --no-stage

# Generate a changelog with multiline commit messages
uplift changelog --multiline
```

## Flags

```text
    --all               generate a changelog from the entire history of this
                        repository
    --diff-only         output the changelog diff only
    --exclude strings   a list of regexes for excluding conventional commits
                        from the changelog
-h, --help              help for changelog
    --include strings   a list of regexes to cherry-pick conventional commits
                        for the changelog
    --multiline         include multiline commit messages within changelog (skips truncation)
    --sort string       the sort order of commits within each changelog entry
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
