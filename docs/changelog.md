# Generating a Changelog

Uplift can generate or amend a changelog (`CHANGELOG.md`) for your repository based on the [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format.

```sh
uplift changelog
```

## Excluding Commits

You can exclude commits from the changelog by specifying a list of commit prefixes.

```sh
uplift changelog --exclude chore,ci,test
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

:octicons-beaker-24: Experimental

If your repository does not contain a `CHANGELOG.md` file, you can generate one that spans its entire history. A word of warning, this does require a tagging structure to be in place.

```sh
uplift changelog --all
```
