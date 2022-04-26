# Git

Uplift requires a git repository to be cloned in a certain way to ensure all features run as expected. However, since only a subset of features may be needed, existing git checks can be suppressed when needed.

## ignoreDetached

```yaml linenums="1"
# .uplift.yml

git:
  # A flag for suppressing the git detached HEAD repository check. If set to
  # true, uplift will report a warning while running, otherwise uplift will
  # raise an error and stop.
  #
  # Defaults to false
  ignoreDetached: true
```

## ignoreShallow

```yaml linenums="1"
# .uplift.yml

git:
  # A flag for suppressing the git shallow repository check. If set to true,
  # uplift will report a warning while running, otherwise uplift will raise
  # an error and stop.
  #
  # Defaults to false
  ignoreShallow: true
```
