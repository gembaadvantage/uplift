# Travis CI

An example YAML file[^1] for configuring uplift to run on [Travis CI](https://www.travis-ci.com/). Access to GitHub is managed through their dedicated [GitHub Application](https://docs.travis-ci.com/user/tutorial/#to-get-started-with-travis-ci-using-github). As uplift requires write permissions to your repository, a [Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) needs to be configured with the `public_repo` permission and added to Travis CI as an [encrypted variable](https://docs.travis-ci.com/user/environment-variables/#defining-encrypted-variables-in-travisyml).

```yaml
# .travis.yml

# Setup the pipeline based on your chosen language
language: go

git:
  depth: false

before_install:
  - curl https://raw.githubusercontent.com/gembaadvantage/uplift/main/scripts/install | bash

script:
  - git remote set-url origin https://${GH_UPLIFT}@github.com/${TRAVIS_REPO_SLUG}.git
  - git checkout $TRAVIS_BRANCH

deploy:
  - provider: script
    skip_cleanup: true
    script: uplift release
    on:
      branch: main
      condition: $TRAVIS_OS_NAME = linux # (1)

env:
  global:
    secure: 0l3pSB3Du+YQuV4Gf0R2PoPlrGnmuQhpEbab4KmgUJu6P4S.... # (2)
```

1. If you have configured Travis CI to use a [build matrix](https://docs.travis-ci.com/user/build-matrix/), a condition like this should be used to ensure uplift is only run once.
2. You will need to download travis in order to encrypt variables. Once downloaded, you must first login `travis --login --pro --github-token=<TRAVIS_TOKEN>` and then generate an encrypted variable with a command similar to `echo GH_UPLIFT=<PERSONAL_ACCESS_TOKEN> | travis encrypt --add --pro`

[^1]: There are many different ways of [installing](../install.md) uplift within a pipeline. Sudo access is needed when installing the binary into a protected path such as `/usr/local/bin`
