# Changing how Uplift works with Git

Git is the core component behind Uplift. So it is only fitting that Uplift supports customisation options when interacting with it.

## Skipping Git Checks

Sometimes Uplift will complain about the state of the current git repository. If you wish to override this behaviour, you can customise how Uplift interacts with git. It is worth pointing out that you cannot make Uplift run against a dirty repository, and you will have to fix that yourself using our handy [FAQ](../faq/gitdirty.md)

### Ignoring a Detached HEAD

Either use the `--ignore-detached` flag:

```sh
uplift release --ignore-detached
```

Or include the following entry in your config file:

```yaml linenums="1"
# .uplift.yml

git:
  ignoreDetached: true
```

### Ignoring a Shallow Clone

Either use the `--ignore-shallow` flag:

```sh
uplift release --ignore-shallow
```

Or include the following entry in your config file:

```yaml linenums="1"
# .uplift.yml

git:
  ignoreShallow: true
```

## Additional Git Push Options

Since Git version 2.10, the ability to pass additional push options (`--push-option`) to the remote has been supported. Some SCMs have used this to support custom behaviour after a push. By including the following entry in your config file, Uplift can use these options independently during a push of staged files and a push of a new tag.

```yaml linenums="1"
# .uplift.yml

git:
  pushOptions:
    - option: ci.skip
      skipBranch: true
      skipTag: false
```

## Prevent Staging of Files

To take ownership of staging and committing files to your repository, you can disable this automatic feature within Uplift using the `--no-stage` flag.

```sh
uplift release --no-stage
```

## Prevent Pushing Changes to the Remote

To take ownership of pushing changes to your repository, you can disable this automatic feature within Uplift with the `--no-push` flag.

```sh
uplift release --no-push
```
