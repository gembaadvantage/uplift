# Changing how Uplift works with Git

Sometimes you will find Uplift will complain about the current git repository, and in most situations, it is for the right reason. If you wish to override this behaviour, you can customise how Uplift interacts with git. It is worth pointing out that you cannot make Uplift run against a dirty repository, and you will have to fix that yourself using our handy [FAQ](../faq/gitdirty.md).

## Ignoring a Detached HEAD

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

## Ignoring a Shallow Clone

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
