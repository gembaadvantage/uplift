# Drone

An example YAML file for configuring Uplift to run on [Drone](https://www.drone.io/)[^1]. Instructions for deploying a self-hosted instance of drone can be found [here](https://docs.drone.io/server/provider/github/). Uplift requires write permissions to your repository. A [Personal Access](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) needs to be configured with the `public_repo` permission and added to Drone as an [encrypted secret](https://docs.drone.io/secret/encrypted/).

```{ .yaml .annotate linenums="1" }
# .drone.yml

kind: pipeline
type: docker
name: default

steps:
  - name: set-remote
    image: docker:git
    environment:
      GITHUB_PAT:
        from_secret: github_pat
    commands:
      - CLONE_URL=${DRONE_GIT_HTTP_URL##https://}
      - git remote set-url origin https://$GITHUB_PAT@$CLONE_URL
    when:
      branch:
        - main
      event:
        - push

  - name: release
    image: gembaadvantage/uplift
    commands:
      - uplift release --fetch-all # (1)
    when:
      branch:
        - main
      event:
        - push
---
kind: secret
name: github_pat
data: VLO71Ad3QSALQjELGKg5U7r92823a9e4vmu7xUw3LJ9xKwZu8X...
```

1. Drone does not retrieve any tags by default. Uplift can retrieve all of the latest tags by providing the `--fetch-all` flag.

[^1]: Drone still has a self-hosted offering, but its SaaS product is now integrated with [Harness CI](https://harness.io/).
