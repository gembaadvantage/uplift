# Git repository has a detached HEAD

Uplift may not run reliably[^1] against a git repository that is checked out with a detached HEAD. A detached HEAD occurs when a checkout is made against a specific commit rather than a branch. Many of the documented CI providers use this strategy to ensure a build runs against a commit that triggered it. While in this state, Uplift cannot push changes back to the `main` branch.

```text
uplift cannot reliably run when the repository is in a detached HEAD state. Some features
will not run as expected. To suppress this error, use the '--ignore-detached' flag, or
set the required config.

For further details visit: https://upliftci.dev/faq/gitdetached
```

To resolve this error, you have the following options:

1. If you are using a documented CI provider, view the example YAML configuration to ensure your repository is in the right state before running Uplift. If your CI provider isn't listed, please consult their documentation. We would appreciate it if you contributed back with your findings.
1. You can suppress the error by either setting the global [`--ignore-detached`](../cli/root.md#-ignore-detached) flag or by disabling it in the uplift [config](../config/git.md#ignoredetached) file.

[^1]: Features such as file bumping and changelog management will be impacted.
