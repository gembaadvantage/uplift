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

### --no-prefix

Strips the default `v` prefix from the next calculated semantic version.

```sh
uplift release --no-prefix
```

!!!tip "Should only need to use this once"

    Once a repository has been tagged with either scheme, e.g. `1.0.0` or `v1.0.0`, uplift will continue using it. This flag is most useful when tagging your repository for the first time.

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

### --exclude

By specifying a list of conventional commit prefixes, certain types of commits can be filtered (excluded) from the generated changelog. Optionally include the scope to narrow the range commits that will be excluded. For example, `chore` will exclude any commit with that prefix, while `chore(build)` only exclude commits with the build scope.

```sh
uplift changelog --exclude chore,ci,docs
```

### --sort

Can be used to change the sort order of commits within each changelog entry. By default entries are sorted in descending (`desc`) order. Latest to oldest commit. Supported values are `desc` and `asc`, or any case variant of these.

```sh
uplift changelog --sort desc
```
