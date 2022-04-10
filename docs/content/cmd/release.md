---
sidebar_position: 5
---

# uplift release

Release the next semantic version of your git repository. A release will automatically bump any files and tag the associated commit with the required semantic version

```sh
uplift release [FLAGS]
```

## Flags

### --check

Checks if the latest commit contains a conventional commit prefix that will trigger a new release. Returns an exit code of `0` if a release would be carried out.

```sh
$ uplift release --check
   • check release
      • retrieved latest commit   message=feat: this is a new feature
      • detected releasable commit increment=minor
```

### --fetch-all

Ensure all tags associated with the git repository are fetched before carrying out the release.

```sh
uplift release --fetch-all
```

### --prerelease

Append a prerelease suffix to the next calculated semantic version and use it as a prerelease tag. Supporting the [Semver 2.0.0](https://semver.org/) specification, additional labels can be provided for both the prerelease and metadata parts:

- 1.0.0`-beta.1`
- 1.1.0`-beta.1+20220312`

```sh
uplift release --prerelease beta.1
```

### --skip-bumps

Skip the bumping of any files configured within the uplift configuration file in your repository.

```sh
uplift release --skip-bumps
```

### --skip-changelog

Skip the creation or updating of a changelog during the release.

```sh
uplift release --skip-changelog
```
