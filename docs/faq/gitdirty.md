# Git Repository is in a Dirty State

Uplift can't run against a git repository with unstaged or uncommitted files, typically known as a dirty state. If detected, Uplift will report the following error:

```text
uplift cannot reliably run if the repository is in a dirty state. Changes detected:
?? {==coverage.out==}

Please check and resolve the status of these files before retrying. For further
details visit: https://upliftci.dev/faq/gitdirty
```

## How to fix it

You can resolve this error in one of two ways.

### Use a .gitignore file

Add or modify an existing `.gitignore` file to ignore the offending files listed in the error message.

### Adapt your CI

- Ensure no tracked files are unexpectedly modified
- Prevent the creation of temporary files. If this isn't possible, you can fall back to using a `.gitignore` file.
