# Bumping your Files

If you only need to bump the semantic version within specific files, Uplift has you covered. A `.uplift.yml` configuration file is required for this to work.

```yaml linenums="1"
# .uplift.yml
# Example of bumping a package.json file

bumps:
  - file: package.json
    json:
      - path: "version"
        semver: true
```

```sh
uplift bump
```

Please review our guide on configuring file bumps for comprehensive details.

## Prerelease Support

:octicons-beaker-24: Experimental

Uplift has early support for bumping files with prerelease metadata. You will need to calculate this upfront.

```sh
uplift bump --prerelease beta.1+20220930
```
