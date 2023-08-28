# Uplift Configuration

You are free to control Uplift through the use of a dedicated configuration file. A variety of different naming conventions are supported:

- `.uplift.yml`
- `.uplift.yaml`
- `uplift.yml`
- `uplift.yaml`

## annotatedTags

```{ .yaml .annotate linenums="1" }
# Use annotated tags instead of lightweight tags when tagging a new
# semantic version. An annotated tag is treated like a regular commit
# by git and contains both author details and a commit message. Uplift
# will either use its defaults or the custom commit details provided
# when generated the annotated tag.
#
# Defaults to false
annotatedTags: true
```

## bumps

```{ .yaml .annotate linenums="1" }
# Define a series of files whose semantic version will be bumped.
# Supports both Regex and JSON Path based file bumps
#
# Defaults to no files being bumped
bumps:
  # The path of the file relative to where Uplift is executed. Glob
  # patterns can be used to match multiple files at the same time
  - file: package.json

    # A JSON path matcher should be used when bumping the file. Multiple
    # path matches are supported. Each will be carried out in the order
    # they are defined here. All matches must succeed for the file to
    # be bumped. JSON path syntax is based on
    # https://github.com/tidwall/sjson
    #
    # Defaults to no matchers
    json:
      # A JSON path that will be used for matching the version that
      # will be replaced within the file
      - path: "version"

        # If the matched version in the file should be replaced with a
        # semantic version. This will strip any 'v' prefix if needed
        #
        # Defaults to false
        semver: true

  # The path of the file relative to where Uplift is executed. Glob
  # patterns can be used to match multiple files at the same time
  - file: chart/my-chart/Chart.yaml

    # A regex matcher should be used when bumping the file. Multiple
    # regex matches are supported. Each will be carried out in the order
    # they are defined here. All matches must succeed for the file to
    # be bumped
    #
    # Defaults to no matchers
    regex:
      # The regex that should be used for matching the version that
      # will be replaced within the file
      - pattern: "version: $VERSION"

        # If the matched version in the file should be replaced with a
        # semantic version. This will strip any 'v' prefix if needed
        #
        # Defaults to false
        semver: true

        # The number of times any matched version should be replaced
        #
        # Defaults to 0, which replaces all matches
        count: 1
```

## changelog

```{ .yaml .annotate linenums="1" }
# Customise the creation of the Changelog
changelog:
  # Change the sort order of the commits within each changelog entry.
  # Supported values are asc or desc (case is ignored)
  #
  # Defaults to desc (descending order) to mirror the default behaviour
  # of "git log"
  sort: asc

  # A list of commits to exclude during the creation of a changelog.
  # Provide a list of regular expressions for matching commits that
  # are to be excluded. Auto-generated commits from Uplift
  # (with the prefix ci(uplift)) will always be excluded
  #
  # Defaults to an empty list. All commits are included
  exclude:
    - '^chore\(deps\)'
    - ^docs
    - ^ci

  # A list of commits to cherry-pick and include during the creation
  # of a changelog. Provide a list of regular expressions for matching
  # commits that are to be included
  include:
    - '^.*\(scope\)'

  # Include multiline commit messages within the changelog. Disables
  # default behaviour of truncating a commit message to its first line
  multiline: true
```

## commitAuthor

```{ .yaml .annotate linenums="1" }
# Changes the commit author used by Uplift when committing any staged
# changes.
#
# Defaults to the Uplift Bot: uplift-bot <uplift@gembaadvantage.com>
commitAuthor:
  # Name of the author
  #
  # Defaults to the author name within the last commit
  name: "joe.bloggs"

  # Email of the author
  #
  # Defaults to the author email within the last commit
  email: "joe.bloggs@gmail.com"
```

## commitMessage

```{ .yaml .annotate linenums="1" }
# Changes the default commit message used by Uplift when committing
# any staged changes.
#
# Default commit message is: ci(uplift): uplifted for version v0.1.0
commitMessage: "chore(release): this is a custom release message"
```

## hooks

```{ .yaml .annotate linenums="1" }
# All hooks default to an empty list and will be skipped
hooks:
  # A list of shell commands or scripts to execute before Uplift runs
  # any tasks within any workflow
  before:
    - cargo fetch
    - ENV=VALUE ./my-custom-script.sh
    - bash path//to//my-custom-script.sh # (1)

  # A list of shell commands or scripts to execute before Uplift bumps
  # any configured file
  beforeBump:
    - echo "Before Bump"

  # A list of shell commands or scripts to execute before Uplift runs
  # its changelog generation task
  beforeChangelog:
    - echo "Before Changelog"

  # A list of shell commands or scripts to execute before Uplift tags
  # the repository with the next semantic release
  beforeTag:
    - echo "Before Tag"

  # A list of shell commands or scripts to execute after Uplift
  # completes all tasks within any workflow
  after:
    - echo "After Workflow"

  # A list of shell commands or scripts to execute after Uplift bumps
  # any configured file
  afterBump:
    - echo "After Bump"

  # A list of shell commands or scripts to execute after Uplift generates
  # a new changelog
  afterChangelog:
    - echo "After Changelog"

  # A list of shell commands or scripts to execute after Uplift tags
  # the repository with the next semantic release
  afterTag:
    - echo "After Tag"
```

1. An example of using POSIX-based windows commands is through the [mvdan/sh](https://github.com/mvdan/sh) GitHub library. Pay special attention to the use of `//` when specifying a path

## env

```{ .yaml .annotate linenums="1" }
# Define a set of environment variables that are made available to all
# hooks. Supports loading environment variables from DotEnv (.env)
# files. Environment variables are merged with system wide ones.
env:
  - VARIABLE=VALUE
  - ANOTHER_VARIABLE=ANOTHER VALUE
  - .env
  - path/to/other.env
```

## git

```{ .yaml .annotate linenums="1" }
# Customise how Uplift interacts with Git
git:
  # A flag for suppressing the git detached HEAD repository check. If set
  # to true, Uplift will report a warning while running, otherwise Uplift
  # will raise an error and stop.
  #
  # Defaults to false
  ignoreDetached: true

  # A flag for suppressing the git shallow repository check. If set to
  # true, Uplift will report a warning while running, otherwise Uplift
  # will raise an error and stop.
  #
  # Defaults to false
  ignoreShallow: true

  # An array of Git push options that can be independently configured
  # for both branch and tag operations within Uplift. Provided options
  # will be filtered accordingly and appended to the git push operation
  # through the use of the --push-option flag as documented in
  # https://git-scm.com/docs/git-push#Documentation/git-push.txt
  #
  # Defaults to any empty array
  pushOptions:
    - ci.skip

    # For more fine-grained control, a push option can be skipped
    # based on the type of push being executed
    - option: ci.variable="MAX_RETRIES=10"
      skipTag: true
      skipBranch: false
```

## gitea

```{ .yaml .annotate linenums="1" }
# Add support for Gitea SCM detection
gitea:
  # The URL of the self-hosted instance of Gitea. Only the scheme and
  # hostname are required. The hostname is used when matching against
  # the configured remote origin of the cloned repository
  #
  # Defaults to empty string i.e. no detection is supported
  url: https://my.gitea.com
```

## github

```{ .yaml .annotate linenums="1" }
# Add support for GitHub SCM detection
github:
  # The URL of the enterprise instance of GitHub. Only the scheme and
  # hostname are required. The hostname is used when matching against
  # the configured remote origin of the cloned repository
  #
  # Defaults to empty string i.e. no detection is supported
  url: https://my.github.com
```

## gitlab

```{ .yaml .annotate linenums="1" }
# Add support for GitLab SCM detection
gitlab:
  # The URL of the self-managed instance of GitLab. Only the scheme and
  # hostname are required. The hostname is used when matching against
  # the configured remote origin of the cloned repository
  #
  # Defaults to empty string i.e. no detection is supported
  url: https://my.gitlab.com
```
