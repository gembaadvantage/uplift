# Uplift Configuration

You are free to control Uplift through the use of a dedicated configuration file. A variety of different naming conventions are supported:

- `.uplift.yml`
- `.uplift.yaml`
- `uplift.yml`
- `uplift.yaml`

```yaml linenums="1"
# .uplift.yml

# Define a set of environment variables that are made available to all
# hooks. Supports loading environment variables from DotEnv (.env)
# files. Environment variables are merged with system wide ones.
env:
  - VARIABLE=VALUE
  - ANOTHER_VARIABLE=ANOTHER VALUE
  - .env
  - path/to/other.env

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

# Use annotated tags instead of lightweight tags when tagging a new
# semantic version. An annotated tag is treated like a regular commit
# by git and contains both author details and a commit message. Uplift
# will either use its defaults or the custom commit details provided
# when generated the annotated tag.
#
# Defaults to false
annotatedTags: true

# Customise how Uplift responds to its inbuilt Git checks
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

# Customise the creation of the Changelog
changelog:
  # Change the sort order of the commits within each changelog entry.
  # Supported values are asc or desc (case is ignored)
  #
  # Defaults to desc (descending order) to mirror the default behaviour
  # of "git log"
  sort: asc

  # A list of commits to exclude during the creation of a changelog.
  # Provide a list of conventional commit prefixes to filter on.
  # Auto-generated commits from Uplift (with the prefix ci(uplift)) will
  # always be excluded
  #
  # Defaults to including all commits within the generated changelog
  exclude:
    - chore(deps)
    - docs
    - ci

# Define a series of files whose semantic version will be bumped.
# Supports both Regex and JSON Path based file bumps
#
# Defaults to no files being bumped
bumps:
  # The path of the file relative to where Uplift is executed
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

  # The path of the file relative to where Uplift is executed
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