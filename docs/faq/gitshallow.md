# Git Repository contains a Shallow Clone

Uplift may not run reliably[^1] against a git repository that contains a shallow clone. A shallow clone will result in a git repository containing history about the latest commit only. Many of the documented CI providers use this strategy to improve the efficiency of a clone, especially for large repositories.

```text
uplift cannot reliably run against a shallow clone of the repository. Some features may not
work as expected. To suppress this error, use the '--ignore-shallow' flag, or set the
required config.

For further details visit: https://upliftci.dev/faq/gitshallow
```

To resolve this error, you have the following options:

1. If you are using a documented CI provider, view the example YAML configuration to ensure your repository is in the right state before running Uplift. If your CI provider isn't listed, please consult their documentation. We would appreciate it if you contributed back with your findings.
2. You can suppress the error by either setting the global [`--ignore-shallow`](../cli/root.md#-ignore-shallow) flag or by disabling it in the Uplift [config](../config/git.md#ignoreshallow) file.

[^1]: Depending on the clone strategy of your CI provider, many, if not all features of Uplift will be impacted.
