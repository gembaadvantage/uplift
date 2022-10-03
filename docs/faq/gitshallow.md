# Git Repository contains a Shallow Clone

A git repository from a shallow clone will contain a truncated commit history and potentially no previous tags, disabling most, if not all, of the Uplift features. Cloning behaviour will differ between CI providers. If detected, Uplift will report the following error:

```text
uplift cannot reliably run against a shallow clone of the repository.
Some features may not work as expected. To suppress this error, use the
'{==--ignore-shallow==}' flag, or set the required {==config==}.

For further details visit: https://upliftci.dev/faq/gitshallow
```

## How to fix it

You can resolve this error in one of three ways.

### Fetch the history

If no history exists, use the `--` flag...

### Fetch the tags

If no tags exist, use the `--fetch-tags` flag to fetch all tags from the origin.

### Suppress the error

You can suppress this error by setting the `--ignore-shallow` flag or by modifying your `.uplift.yml` config file:

```yaml linenums="1"
# .uplift.yml

git:
  ignoreShallow: true
```
