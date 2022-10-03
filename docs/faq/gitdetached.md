# Git repository has a detached HEAD

File bumping and changelog creation will not run reliably against a git repository cloned at a specific commit rather than a branch, known as a detached HEAD. Some CI providers use this as an efficient cloning strategy, but it prevents Uplift from pushing changes back to the default branch. If detected, Uplift will report the following error:

```text
uplift cannot reliably run when the repository is in a detached HEAD state.
Some features will not run as expected. To suppress this error, use the
'{==--ignore-detached==}' flag, or set the required {==config==}.

For further details visit: https://upliftci.dev/faq/gitdetached
```

## How to fix it

You can resolve this error in one of two ways.

### Reattach the HEAD of your Repository

Resolving a detached HEAD requires you to check out the default branch, effectively reattaching the HEAD. Please look at our documented CI providers for examples of how to do this.

### Suppress the error

You can suppress this error by setting the `--ignore-detached` flag or by modifying your `.uplift.yml` config file:

```yaml linenums="1"
# .uplift.yml

git:
  ignoreDetached: true
```
