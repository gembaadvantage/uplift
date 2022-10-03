# Gitea

As [Gitea](https://gitea.io/en-us/) is a self-hosted SCM, custom configuration is required to support detection.

```yaml linenums="1"
# uplift.yml

# Add support for Gitea SCM detection
gitea:
  # The URL of the self-hosted instance of Gitea. Only the scheme and
  # hostname are required. The hostname is used when matching against
  # the configured remote origin of the cloned repository
  #
  # Defaults to empty string i.e. no detection is supported
  url: https://my.gitea.com
```
