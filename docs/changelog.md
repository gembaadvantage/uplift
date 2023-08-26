# Generating a Changelog

Uplift can generate or amend your repository's changelog (`CHANGELOG.md`) based on the [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format.

```sh
uplift changelog
```

## Excluding Commits

You can exclude commits from the changelog by specifying a list of regex. Matching against a commit prefix is the most straightforward approach to doing this.

```sh
uplift changelog --exclude "^chore,^ci,^test"
```

## Including Commits

The inverse behaviour is also supported. You can cherry-pick commits by specifying a list of regex. Helpful if you want to generate a changelog for a particular scope of commits.

```sh
uplift changelog --include "^.*\(scope\)"
```

## Changing the Commit Order

Commits are written to a changelog in descending order, reflecting the behaviour of `git log`. Change this order by specifying `asc` (case insensitive).

```sh
uplift changelog --sort asc
```

## Output the Changelog Diff

You can output the calculated changelog difference (_diff_) to `stdout` without modifying the local repository.

```sh
uplift changelog --diff-only
```

## Migrate an Existing Repository

If your repository does not contain a `CHANGELOG.md` file, you can generate one that spans its entire history. A word of warning, this does require a tagging structure to be in place.

```sh
uplift changelog --all
```

## Supporting Multiline Commits

You can configure `uplift` to include multiline commit messages within your changelog, by disabling its default behaviour to truncate them to a single line.

```sh
uplift changelog --multiline
```
