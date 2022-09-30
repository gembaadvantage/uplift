# Environment Variables

Define environment variables that will be made available to all hooks. Environment variables can be individually listed or defined within [dotenv](https://hexdocs.pm/dotenvy/dotenv-file-format.html)[^1] (.env) files. Uplift will merge all environment variables with any pre-existing system ones.

```yaml linenums="1"
# .uplift.yml

env:
  - VARIABLE=VALUE
  - ANOTHER_VARIABLE=ANOTHER VALUE
  - .env
  - path/to/other.env
```

[^1]: Dotenv support is provided through the [github.com/joho/godotenv](https://github.com/joho/godotenv) library
