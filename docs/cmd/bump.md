---
sidebar_position: 3
---

# uplift bump

Bumps the semantic version within files in your git repository. The version bump is based on the conventional commit message from the last commit. Uplift can bump the version in any file using regex pattern matching

```sh
uplift bump [FLAGS]
```

## Flags

### --prerelease

Append a prerelease suffix to the next calculated semantic version and use it as a prerelease tag. Supporting the [Semver 2.0.0](https://semver.org/) specification, additional labels can be provided for both the prerelease and metadata parts:

- 1.0.0`-beta.1`
- 1.1.0`-beta.1+20220312`

```sh
uplift bump --prerelease beta.1
```
