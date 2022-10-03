# Command Line

Release the next semantic version of your git repository. A release will automatically bump any files and tag the associated commit with the required semantic version

## Usage

```text
uplift release [flags]
```

## Flags

```text
    --check               check if a release will be triggered
    --exclude strings     a list of conventional commit prefixes to exclude from the changelog
    --fetch-all           fetch all tags from the remote repository
-h, --help                help for release
    --no-prefix           strip the default 'v' prefix from the next calculated semantic version
    --prerelease string   append a prerelease suffix to next calculated semantic version
    --skip-bumps          skips the bumping of any files
    --skip-changelog      skips the creation or amendment of a changelog
    --sort string         the sort order of commits within each changelog entry
```

## Global Flags

```text
--config-dir string   a custom path to a directory containing uplift config (default ".")
--debug               show me everything that happens
--dry-run             run without making any changes
--ignore-detached     ignore reported git detached HEAD error
--ignore-shallow      ignore reported git shallow clone error
--no-push             no changes will be pushed to the git remote
--silent              silence all logging
```
