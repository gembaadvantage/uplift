# Extending Uplift with Hooks

Uplift can be extended through the use of hooks. A hook is a specific point during a workflow where Uplift executes adhoc shell commands and scripts. If you need to print the output from any command or script, use the `--debug` flag.

- `before`: a hook that executes before any tasks within the workflow
- `after`: a hook that executes after completing all workflow tasks
- `beforeBump`: a hook that executes before bumping any configured files
- `afterBump`: a hook that executes after bumping all configured files
- `beforeChangelog`: a hook that executes before generating a changelog
- `afterChangelog`: a hook that executes after changelog generation
- `beforeTag`: a hook that executes before tagging the repository
- `afterTag`: a hook that executes after the repository is tagged

```{ .yaml .annotate linenums="1" }
# .uplift.yml

hooks:
  before:
    - cargo fetch
    - ENV=VALUE ./my-custom-script.sh
    - bash path//to//my-custom-script.sh # (1)
```

1. An example of invoking a script using a POSIX-based Windows shell. Pay special attention to the use of `//` when specifying a path

❤️ to the [github.com/mvdan/sh](https://github.com/mvdan/sh) library.

## Injecting Environment Variables

Extend hook support by defining environment variables that Uplift will inject into the runtime environment. Either list environment variables individually or import them through [dotenv](https://hexdocs.pm/dotenvy/dotenv-file-format.html) (.env) files. Uplift will merge all environment variables with any pre-existing system ones.

```yaml linenums="1"
# .uplift.yml

env:
  - VARIABLE=VALUE
  - ANOTHER_VARIABLE=ANOTHER VALUE
  - .env
  - path/to/other.env
```

❤️ to the [github.com/joho/godotenv](https://github.com/joho/godotenv) library.
