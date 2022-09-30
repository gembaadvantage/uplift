# Git Repository is in a Dirty State

Uplift won't run against a git repository that has un-staged and/or un-committed files, typically know as a dirty state. Uplift requires a clean git working directory.

```sh
uplift cannot reliably run if the repository is in a dirty state. Changes detected:
 M main.go
?? coverage.out

Please check and resolve the status of these files before retrying. For further
details visit: https://upliftci.dev/faq/gitdirty
```

As you can see the error message shows the offending files and their current git [status](https://git-scm.com/docs/git-status#_short_format). To resolve the error, you have the following options:

1. Add a `.gitignore` file to your repository to ensure these files are no longer tracked
2. Change your CI approach to ensure no tracked files are modified or temporary files are generated before uplift is run
