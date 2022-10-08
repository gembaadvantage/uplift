# Command Line

Bumps the semantic version within files in your git repository. The version bump is based on the conventional commit message from the last commit. Uplift can bump the version in any file using regex pattern matching

## Usage

```text
uplift bump [flags]
```

## Flags

```text
-h, --help                help for bump
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
