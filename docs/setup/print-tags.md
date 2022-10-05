# Printing Repository Tags

Uplift provides utility functions for printing the current or next calculated semantic version of your repository to `stdout`. Useful if you want to use Uplift alongside other tools in your CI.

## Printing the Next Tag

Scans all commits from the latest release[^1] and prints the next calculated semantic version to `stdout`. Prints an empty string if no commits exist that triggers the next semantic version.

```sh
NEXT_TAG=$(uplift tag --next --silent)
```

## Printing the Current Tag

Scans and prints the most recent semantic version from a repository with existing git tags to `stdout`. Prints an empty string if no tags exist.

```sh
CURRENT_TAG=$(uplift tag --current)
```

## Printing the Tag Transition

Scans all commits from the latest release[^1] and prints the current and next calculated semantic version to `stdout`. Both tags are separated by one whitespace ensuring compatibility with many Linux text processing tools, e.g. `v0.1.0 v0.2.0`. Prints an empty string if no commits exist that triggers the next semantic version.

```sh
TAG_TRANSITION=$(uplift tag --current --next --silent)
```

[^1]: If this is a repository without any previous releases, Uplift will scan the entire commit history
