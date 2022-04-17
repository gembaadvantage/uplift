# GitHub Action

The official [GitHub Action](https://github.com/gembaadvantage/uplift-action) can be used to configure uplift within your workflow. As uplift is designed to push changes back to your repository, you will need to provide it with an access token. This is by [design](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#using-the-github_token-in-a-workflow).

```{ .yaml .annotate linenums="1" }
# .github/workflows/ci.yml

name: ci
on:
  push:
    branches:
      - main
  pull_request:
jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0 # (1)
      - name: Release
        if: github.ref == 'refs/heads/main'
        uses: gembaadvantage/uplift-action@v2.0.1
        with:
          args: release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }} # (2)
```

1. Setting a `fetch-depth` of 0 will ensure all tags are retrieved which is required by uplift to determine the next semantic version
2. When you use the repository's `GITHUB_TOKEN` to perform tasks, events triggered by the `GITHUB_TOKEN` will not create a new workflow run.

## Triggering another Workflow

To ensure uplift triggers another workflow run when tagging the repository, a [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) should be created and stored as a [secret](https://docs.github.com/en/actions/security-guides/encrypted-secrets). This will then replace the default `GITHUB_TOKEN` as follows:

```{ .yaml .annotate linenums="1" hl_lines="23" }
# .github/workflows/ci.yml

name: ci
on:
  push:
    branches:
      - main
  pull_request:
jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Release
        if: github.ref == 'refs/heads/main'
        uses: gembaadvantage/uplift-action@v2.0.1
        with:
          args: release
        env:
          GITHUB_TOKEN: ${{ secrets.GH_UPLIFT }}
```
