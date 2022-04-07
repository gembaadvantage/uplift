---
sidebar_position: 4
---

# uplift changelog

Create or update an existing changelog with an entry for the latest semantic release. For a first release, all commits between the latest tag and trunk will be written to the changelog. Subsequent entries will contain only commits between release tags.

```sh
uplift changelog [FLAGS]
```

## Flags

### --all

Generates a changelog from the entire history of the git repository. Provides a great way to migrate an existing changelog process to uplift.

```sh
uplift changelog --all
```

### --diff-only

Writes the calculated changelog diff to stdout without modifying the current repository. Useful for combining uplift with any other tooling.

```sh
$ uplift changelog --diff-only
   • changelog
      • determine changes for release tag=1.0.0
      • changeset identified      commits=3 date=2022-03-25 tag=1.0.0
## [1.0.0] - 2022-03-25

- `e988091` feat: a brand new feature
- `11d039b` ci: speed up existing workflow
- `fad2c38` docs: update to existing documentation
```

Easily capture the diff within an environment variable for use within a custom script:

```sh
CHANGELOG_DIFF=$(uplift changelog --diff-only --silent)
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

## SCM Detection

During changelog creation uplift will attempt to identify the SCM provider associated with the repository. Upon successful detection, uplift will embed links to tags and commits within the changelog, making it easier to inspect them from within the SCM tool.

Supported SCM providers:

- GitHub (Cloud)
- GitLab (Cloud)
- AWS CodeCommit
