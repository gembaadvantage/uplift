# Command Line

```text
Release the next semantic version of your git repository. A release consists of
a three-stage process. First, all configured files will be bumped (patched)
using the next semantic version. Second, a changelog entry containing all
commits for the latest semantic release will be created. Finally, Uplift will
tag the repository. Uplift automatically handles the staging and pushing of
modified files and the tagging of the repository with two separate git pushes.
But this behavior can be disabled to manage these actions manually.

Parts of this release process can be disabled if needed.

https://upliftci.dev/first-release/
```

## Usage

```text
uplift release [flags]
```

## Examples

```text
# Release the next semantic version
uplift release

# Release the next semantic version without bumping any files
uplift release --skip-bumps

# Release the next semantic version without generating a changelog
uplift release --skip-changelog

# Append a prerelease suffix to the next calculated semantic version
uplift release --prerelease beta.1

# Ensure any "v" prefix is stripped from the next calculated semantic
# version to explicitly adhere to the SemVer specification
uplift release --no-prefix
```

## Flags

```text
    --check                       check if a release will be triggered
    --exclude strings             a list of regexes for excluding conventional
                                  commits from the changelog
    --fetch-all                   fetch all tags from the remote repository
-h, --help                        help for release
    --include strings             a list of regexes to cherry-pick conventional
                                  commits for the changelog
    --multiline                   include multiline commit messages within
                                  changelog (skips truncation)
    --no-prefix                   strip the default 'v' prefix from the next
                                  calculated semantic version
    --prerelease string           append a prerelease suffix to next calculated
                                  semantic version
    --skip-bumps                  skips the bumping of any files
    --skip-changelog              skips the creation or amendment of a changelog
    --skip-changelog-prerelease   skips the creation of a changelog entry for a
                                  prerelease
    --sort string                 the sort order of commits within each
                                  changelog entry
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
