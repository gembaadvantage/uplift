# Command Line

```text
Semantic versioning the easy way.
```

## Usage

```text
uplift [command]
```

## Flags

```text
    --config-dir string            a custom path to a directory containing
                                   uplift config (default ".")
    --debug                        show me everything that happens
    --dry-run                      run without making any changes
-h, --help                         help for uplift
    --ignore-detached              ignore reported git detached HEAD error
    --ignore-existing-prerelease   ignore any existing prerelease when
                                   calculating next semantic version
    --ignore-shallow               ignore reported git shallow clone error
    --no-push                      no changes will be pushed to the git remote
    --no-stage                     no changes will be git staged
    --silent                       silence all logging
```

## Commands

```text
bump        Bump the semantic version within files
changelog   Create or update a changelog with the latest semantic release
completion  Generate completion script for your target shell
help        Help about any command
release     Release the next semantic version of a repository
tag         Tag a git repository with the next semantic version
version     Prints the build time version information
```
