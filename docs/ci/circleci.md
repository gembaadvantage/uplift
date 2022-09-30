# CircleCI

An example YAML file[^1] for configuring uplift to run on [CircleCI](https://circleci.com). As uplift is designed to push changes back to your GitHub repository you will need to ensure CircleCI is granted [write access](https://circleci.com/docs/2.0/gh-bb-integration) to your repository.

```{ .yaml .annotate linenums="1" }
# .circleci/config.yml

version: 2.1
workflows:
  main:
    jobs:
      - release:
          filters:
            branches:
              # Only trigger on the main branch
              only: main
jobs:
  release:
    docker:
      # Can use whatever base image you like
      - image: cimg/go:1.18
    steps:
      # Configure an SSH key that provides write access to your GitHub repository
      - add_ssh_keys:
          fingerprints:
            - "3b:c7:44:c9:34:ab:a4:fd:6c:33:4e:a7:7a:97:79:55" # (1)
      - checkout
      # Additional actions specified here
      - run: curl https://raw.githubusercontent.com/gembaadvantage/uplift/main/scripts/install | bash
      - run: uplift release
```

1. By default CircleCI will only have read-only access to your repository. For uplift to work, write access is required. This can be achieved by accessing a repository as a [machine-user](https://circleci.com/docs/2.0/gh-bb-integration/#controlling-access-via-a-machine-user) and then loading its [SSH key](https://circleci.com/docs/2.0/configuration-reference/#add-ssh-keys) into the pipeline by its fingerprint

[^1]: There are many different ways of [installing](../install.md) uplift within a pipeline. Sudo access is needed when installing the binary into a protected path such as `/usr/local/bin`
