---
sidebar_position: 2
---

# Basics

Ideally none of these values should need to be set as uplift will always use what we believe to be sensible default values.

## firstVersion

```yaml
# .uplift.yml

# An initial version that will be used as the first tag within your repository.
# Tags with a 'v' prefix are supported and will be treated as a semantic version
# by uplift. For existing repositories with a semantic version scheme already in
# place this setting will be ignored.
#
# Defaults to 0.1.0
firstVersion: v1.0.0
```

## commitMessage

```yaml
# .uplift.yml

# Changes the default commit message used by uplift when committing any staged
# changes. If carrying out a full release, uplift will have staged and committed
# any changelog creation or amendments, and any file within the repository that
# has been bumped.
#
# Defaults to ci(uplift): uplifted for version <LATEST_TAG>
commitMessage: "chore: a custom commit message"
```

## commitAuthor

```yaml
# .uplift.yml

# Changes the commit author used by uplift when committing any staged changes.
# Useful if you want uplift to be treated as a bot and own its commits.
#
# Defaults to the author who last committed to the repository and effectively
# triggered uplift. Uplift likes to give credit to the author that is releasing
# a new feature or fixing that bug they found.
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

## annotatedTags

```yaml
# .uplift.yml

# Use annotated tags instead of lightweight tags when tagging a new semantic
# version. An annotated tag is treated like a regular commit by git and contains
# both author details and a commit message. Uplift will either use its defaults
# or the custom commit details provided when generated the annotated tag.
#
# Defaults to false
annotatedTags: true
```

!!!info "What are Annotated Tags?"

    To find out more about annotated tags I recommend reading the official [Git](https://git-scm.com/book/en/v2/Git-Basics-Tagging) documentation
