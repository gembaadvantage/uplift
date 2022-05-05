# GitLab

Uplift comes with built-in detection for GitLab (SaaS). However, when using [Self-Managed GitLab](https://about.gitlab.com/install/), custom configuration is needed.

```yaml linenums="1"
# uplift.yml

gitlab:
  # The URL of the self-hosted instance of GitLab. Only the scheme and
  # hostname are required. The hostname is used when matching against
  # the configured remote origin of the cloned repository
  #
  # Defaults to empty string i.e. no detection is supported
  url: https://my.gitlab.com
```
