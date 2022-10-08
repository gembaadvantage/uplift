# Command Line

Tags a git repository with the next semantic version. The next tag is calculated using the conventional commit message from the last commit.

## Usage

```text
uplift tag [flags]
```

## Flags

```text
   --fetch-all           fetch all tags from the remote repository
-h, --help                help for tag
   --next                output the next tag only
   --no-prefix           strip the default 'v' prefix from the next calculated semantic version
   --prerelease string   append a prerelease suffix to next calculated semantic version
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
