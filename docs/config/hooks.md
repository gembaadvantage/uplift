# Hooks

Provides a mechanism to extend the functionality of uplift through adhoc shell commands and scripts[^1]. Only the `before` hook precedes the git checks within uplift. All temporary files must therefore be ignored using a `.gitignore` file, otherwise the repository will be deemed in a [dirty state](../faq/gitdirty.md) and the release will stop. Hooks tied to either of the `bump`, `changelog` or `tag` operations will be skipped along with its counterpart when needed by uplift.

```{ .yaml .annotate linenums="1" }
# .uplift.yml

# All hooks default to an empty list and will be skipped
hooks:
  # A list of shell commands or scripts to execute before uplift runs
  # any tasks within its release workflow
  before:
    - npm install
    - go mod tidy
    - cargo fetch
    - ENV=VALUE ./my-custom-script.sh
    - bash path//to//my-custom-script.sh # (1)

  # A list of shell commands or scripts to execute before uplift bumps
  # any configured file
  beforeBump:
    - ...

  # A list of shell commands or scripts to execute before uplift runs
  # its changelog generation task
  beforeChangelog:
    - ...

  # A list of shell commands or scripts to execute before uplift tags
  # the repository with the next semantic release
  beforeTag:

  # A list of shell commands or scripts to execute after uplift
  # completes all tasks within its release workflow
  after:
    - ...

  # A list of shell commands or scripts to execute after uplift bumps
  # any configured file
  afterBump:
    - ...

  # A list of shell commands or scripts to execute after uplift generates
  # a new changelog
  afterChangelog:
    - ...

  # A list of shell commands or scripts to execute after uplift tags
  # the repository with the next semantic release
  afterTag:
    - ...
```

1. An example of using POSIX based windows commands through the [mvdan/sh](https://github.com/mvdan/sh) GitHub library. Pay special attention to the use of `//` when specifying a path

!!!tip "Need extra output?"

    Use the `--debug` flag to print output from any of the executed shell commands or scripts

[^1]: Interpretation and execution of shell commands and scripts is carried out through the [mvdan/sh](https://github.com/mvdan/sh) GitHub library.
