# GitLab

An example YAML file for configuring Uplift to run on [GitLab](https://gitlab.com/). To ensure Uplift can push changes back to your repository, you will need to provide it with a [project](https://docs.gitlab.com/ee/user/project/settings/project_access_tokens.html)](https://docs.gitlab.com/ee/user/project/settings/project_access_tokens.html) or [group](https://docs.gitlab.com/ee/user/group/settings/group_access_tokens.html) access token[^1] with the `write_repository` permission.

```{ .yaml .annotate linenums="1" hl_lines="11" }
# .gitlab-ci.yml

stages:
  - release

release:
  stage: release
  image:
    name: gembaadvantage/uplift
    entrypoint: [""]
  dependencies: [] # (1)
  before_script:
    - PROJECT_URL=${CI_PROJECT_URL#"https://"}
    - git remote set-url origin "https://oauth2:${GL_UPLIFT}@${PROJECT_URL}.git"
  variables:
    # Disable shallow cloning of repository
    GIT_DEPTH: 0
  script:
    # GitLab by default checks out a detached HEAD
    - git checkout $CI_COMMIT_REF_NAME
    - uplift release
  # Only run on the default branch of the repository
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event"'
      when: never
    - if: "$CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH"
      when: on_success
```

1. Prevents any dependencies, such as reports, from being unnecessarily copied into this job and causing the git checks to fail

To expose your access token within your pipeline you should add a CI/CD [variable](https://docs.gitlab.com/ee/ci/variables/). In the above example, the access token is exposed through the `GL_UPLIFT` variable.

[^1]: It is best security practice to create an access token with the shortest possible expiration date.
