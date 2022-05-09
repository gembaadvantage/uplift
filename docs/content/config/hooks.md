# Hooks

Provides a mechanism to extend the functionality of uplift through adhoc shell commands and scripts[^1]. Any temporary files must be ignored using a `.gitignore` file, otherwise uplift will deem the repository is in a [dirty state](../faq/gitdirty.md) and stop the release.

```yaml linenums="1"
# .uplift.yml

hooks:
  # A list of shell commands or scripts to execute before uplift runs
  # any tasks within its release workflow
  #
  # Defaults to an empty list. Hooks will be skipped
  before:
    - npm install
    - go mod tidy
    - cargo fetch
    - ENV=VALUE ./my-custom-script.sh
```

!!!tip "Need extra output?"

    Use the `--debug` flag to print output from any of the executed shell commands or scripts

[^1]: Interpretation and execution of shell commands and scripts is carried out through the [mvdan/sh](https://github.com/mvdan/sh) GitHub library.
