# Codefresh

An example YAML file for configuring Uplift to run on [Codefresh](https://g.codefresh.io/welcome). To ensure Uplift can push changes back to your repository, you will need to store your Personal Access Token as a [shared configuration](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/shared-configuration/) and [expose](https://codefresh.io/docs/docs/configure-ci-cd-pipeline/shared-configuration/#using-shared-environment-variables) it to your pipeline as an environment variable, which in this example is `GH_UPLIFT`.

```{ .yaml .annotate linenums="1" }
# codefresh.yml

version: "1.0"
stages:
  - prepare
  - release
steps:
  main_clone: # (1)
    title: "Checkout"
    type: git-clone
    repo: "${{CF_REPO_OWNER}}/${{CF_REPO_NAME}}"
    revision: "${{CF_BRANCH}}" # (2)
    stage: prepare
  uplift:
    title: "Release"
    stage: release
    image: "gembaadvantage/uplift"
    commands:
      - REMOTE_URL=$(git config --get remote.origin.url)
      - CLONE_URL=${REMOTE_URL#"https://"}
      - git remote set-url origin "https://${GH_UPLIFT}@${CLONE_URL}"
      - uplift release
```

1. `main_clone` is a [reserved](https://codefresh.io/docs/docs/codefresh-yaml/steps/git-clone/#basic-clone-step-project-based-pipeline) step within Codefresh and is used to simplify the checkout process. A custom checkout can be performed, but you will need to managed the [working directory](https://codefresh.io/docs/docs/yaml-examples/examples/git-checkout-custom/) across all other steps.
2. By ensuring the clone is of a specific branch, it prevents the checkout of a detached HEAD.
