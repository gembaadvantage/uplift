# Bumping your Files

If you only need to bump the semantic version within specific files, Uplift has you covered. A `.uplift.yml` configuration file is required for this to work. Bumping files using JSON Paths and Regex are currently supported.

```yaml linenums="1"
# .uplift.yml

bumps:
  - file: package.json
    json:
      - path: "version"
        semver: true

  - file: chart/my-chart/Chart.yaml
    regex:
      - pattern: "version: $VERSION"
        semver: true
        count: 1
```

```sh
uplift bump
```

Please review our comprehensive [guide](./reference/config.md#bumps) on configuring file bumps.

❤️ to the [github.com/tidwall/sjson](https://github.com/tidwall/sjson) library.

## Glob Support

If you need to bump multiple similar files at the same time, you can specify a file path using a Glob pattern.

```yaml linenums="1"
# .uplift.yml

bumps:
  - file: "**/package.json"
    json:
      - path: "version"
        semver: true
```

❤️ to the [github.com/goreleaser/fileglob](https://github.com/goreleaser/fileglob) library.

## The $VERSION Token

Writing a regex can be challenging at most times, so Uplift provides the `$VERSION` token for matching a semantic version with an optional `v` prefix. You can include this in any pattern you define within your config.

## Prerelease Support

:octicons-beaker-24: Experimental

Uplift has early support for bumping files with prerelease metadata. You will need to calculate this upfront.

```sh
uplift bump --prerelease beta.1+20220930
```
