---
sidebar_position: 4
---

# Changelog

Uplift can generate and maintain a changelog within your repository for you. We believe that in most situations, these values will never need to be set as uplift uses sensible defaults.

## sort

```yaml linenums="1"
# .uplift.yml

changelog:
  # Change the sort order of the commits within each changelog entry. Supported
  # values are asc or desc (case is ignored)
  #
  # Defaults to desc (descending order) to mirror the default behaviour of "git log"
  sort: asc
```

## exclude

```yaml linenums="1"
# .uplift.yml

changelog:
  # A list of commits to exclude during the creation of a changelog. Provide a list
  # of conventional commit prefixes to filter on. Auto-generated commits from uplift
  # (with the prefix ci(uplift)) will always be excluded
  #
  # Defaults to including all commits within the generated changelog
  exclude:
    - chore(deps)
    - docs
    - ci
```
