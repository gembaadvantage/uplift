# Tagging your Repository

If you only need to manage the tags on your repository, Uplift has you covered.

```sh
uplift tag
```

If you don't want the `v` prefix, no problem; remove it by using the `--strip-prefix` flag.

## Annotated Tags

:octicons-beaker-24: Experimental

Lightweight tags are created by default, equivalent to running `git tag` against your repository. If you need something heavier, such as an annotated tag, modify your configuration file:

```yaml linenums="1"
# .uplift.yml

annotatedTags: true
```

## Prerelease Support

:octicons-beaker-24: Experimental

Uplift has early support for tagging a repository with prerelease metadata. You will need to calculate this upfront.

```sh
uplift tag --prerelease beta.1+20220930
```

If you need Uplift to ignore any existing prerelease metadata when calculating the next semantic version, you must include the `--ignore-existing-prerelease` flag:

```sh
uplift tag --prerelease beta.1+20221006 --ignore-existing-prerelease
```
