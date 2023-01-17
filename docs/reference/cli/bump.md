# Command Line

```text
Calculates the next semantic version based on the conventional commits since the
last release (or identifiable tag) and bumps (or patches) a configurable set of
files with said version. JSON Path or Regex Pattern matching is supported when
scanning files for an existing semantic version. Uplift automatically handles
the staging and pushing of modified files to the git remote, but this behavior
can be disabled, to manage this action manually.

Configuring a bump requires an Uplift configuration file to exist within the
root of your project:

https://upliftci.dev/bumping-files/
```

## Usage

```text
uplift bump [flags]
```

## Examples

```text
# Bump (patch) all configured files with the next calculated semantic version
uplift bump

# Append a prerelease suffix to the next calculated semantic version
uplift bump --prerelease beta.1

# Bump (patch) all configured files but do not stage or push any changes
# back to the git remote
uplift bump --no-stage
```

## Flags

```text
-h, --help                help for bump
    --prerelease string   append a prerelease suffix to next calculated
                          semantic version
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
