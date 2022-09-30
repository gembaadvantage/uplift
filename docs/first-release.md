# Creating your First Release

A release comprises three stages:

1. Patching the semantic version within a set of configured files (_known as file bumping_)
1. Generating a changelog
1. Tagging the repository

## Uplift Configuration

File bumping currently requires a configuration file named `.uplift.yml`. Please review our guide on configuring file bumps for comprehensive details.

```yaml
# .uplift.yml
# Example of bumping a package.json file

bumps:
  - file: package.json
    json:
      - path: "version"
        semver: true
```

Go, create that release 🚀

```sh
uplift release
```

## Skipping Stages

You can skip file bumping `--skip-bumps` and changelog creation `--skip-changelog` by using either of the supported flags.
