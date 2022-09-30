# GitHub

Uplift comes with built-in detection for GitHub (SaaS). However, when using [GitHub Enterprise](https://github.com/enterprise), custom configuration is needed.

```yaml linenums="1"
# uplift.yml

github:
  # The URL of the enterprise instance of GitHub. Only the scheme and
  # hostname are required. The hostname is used when matching against
  # the configured remote origin of the cloned repository
  #
  # Defaults to empty string i.e. no detection is supported
  url: https://my.github.com
```
