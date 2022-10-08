# Command Line

Create or update an existing changelog with an entry for the latest semantic release. For a first release, all commits between the latest tag and trunk will be written to the changelog. Subsequent entries will contain only commits between release tags.

## Usage

```text
uplift changelog [flags]
```

## Flags

```text
   --all               generate a changelog from the entire history of this repository
   --diff-only         output the changelog diff only
   --exclude strings   a list of conventional commit prefixes to exclude
-h, --help              help for changelog
   --sort string       the sort order of commits within each changelog entry
```

## Global Flags

```text
--config-dir string            a custom path to a directory containing uplift config (default ".")
--debug                        show me everything that happens
--dry-run                      run without making any changes
--ignore-detached              ignore reported git detached HEAD error
--ignore-existing-prerelease   ignore any existing prerelease when calculating next semantic version
--ignore-shallow               ignore reported git shallow clone error
--no-push                      no changes will be pushed to the git remote
--no-stage                     no changes will be git staged
--silent                       silence all logging
```
