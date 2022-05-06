# Hooks

Provides a mechanism to extend the functionality of uplift through adhoc shell commands and scripts. Any temporary files must be ignored using a `.gitignore` file, otherwise uplift will deem the repository is in a [dirty state](../faq/gitdirty.md) and stop the release.

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
    - ./my-custom-script.sh
```
