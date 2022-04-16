---
sidebar_position: 3
---

# File Bumping

Uplift can bump the semantic version within any file in your repository using the currently supported schemes. Even though these configs are shown separately, you are free to mix and match within your uplift configuration file.
## regex

```yaml linenums="1"
# .uplift.yml

bumps:
  - # The path of the file relative to where uplift is executed
    file: ./chart/my-chart/Chart.yaml

    # A regex matcher should be used when bumping the file. Multiple regex
    # matches are supported. Each will be carried out in the order they are
    # defined here. All matches must succeed for the file to be bumped
    #
    # Defaults to no matchers
    regex:
      - # The regex that should be used for matching the version that
        # will be replaced within the file
        pattern: "version: $VERSION"

        # If the matched version in the file should be replaced with a semantic version.
        # This will strip any 'v' prefix if needed
        #
        # Defaults to false
        semver: true

        # The number of times any matched version should be replaced
        #
        # Defaults to 0, which replaces all matches
        count: 1
```

### $VERSION

**`$VERSION`** is a placeholder and will match any semantic version, including a version with an optional `v` prefix.

## json

```yaml linenums="1"
# .uplift.yml

bumps:
  - # The path of the file relative to where uplift is executed
    file: ./package.json

    # A JSON path matcher should be used when bumping the file. Multiple path
    # matches are supported. Each will be carried out in the order they are
    # defined here. All matches must succeed for the file to be bumped.
    # JSON path syntax is based on https://github.com/tidwall/sjson
    #
    # Defaults to no matchers
    json:
      - # A JSON path that will be used for matching the version that
        # will be replaced within the file
        path: "version"

        # If the matched version in the file should be replaced with a semantic version.
        # This will strip any 'v' prefix if needed
        #
        # Defaults to false
        semver: true
```

!!!tip "Need more complicated JSON Paths?"

    Uplift uses [SJSON](https://github.com/tidwall/sjson) for setting values through JSON paths. If you need to write more complex JSON paths, don't forget to look at their documentation
