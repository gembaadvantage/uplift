# Semaphore

Example YAML files for configuring uplift to run on [Semaphore 2.0](https://semaphoreci.com/). All Semaphore pipelines start with the default file `.semaphore/semaphore.yml` within your repository. To ensure uplift is only executed on the `main` branch, a separate pipeline YAML file is used, triggered by semaphore [promotions](https://docs.semaphoreci.com/reference/pipeline-yaml-reference/#promotions).

```{ .yaml .annotate linenums="1" hl_lines="19-23" }
# .semaphore/semaphore.yml

version: v1.0
name: CI Pipeline
agent:
  machine:
    type: e1-standard-2
    os_image: ubuntu2004
blocks:
  - name: "CI"
    task:
      jobs:
        - name: "Checkout"
          commands:
            - checkout
        # Additional jobs specified here

# Promotions are used to optionally trigger uplift on any push to the main branch
promotions:
  - name: Uplift
    pipeline_file: uplift.yml
    auto_promote_on:
      - result: passed
        branch:
          - main
```

A dedicated pipeline installs[^1] uplift and triggers a release:

```{ .yaml .annotate linenums="1" }
# .semaphore/uplift.yml

version: "v1.0"
name: Uplift
agent:
  machine:
    type: e1-standard-2
    os_image: ubuntu2004
blocks:
  - name: "Release"
    task:
      prologue:
        commands:
          - checkout
      jobs:
        - name: uplift
          commands:
            - |
              curl https://raw.githubusercontent.com/gembaadvantage/uplift/master/scripts/install > install
              chmod 700 install
              sudo ./install
            - uplift release # (1)
```

1. By default Semaphore installs a GitHub application that has write access to a list of preselected repositories. This ensures no additional configuration is needed to grant uplift permissions for pushing changes back to GitHub

[^1]: There are many different ways of [installing](../install.md) uplift within a pipeline. Sudo access is needed when installing the binary into a protected path such as `/usr/local/bin`
