# Basics

Ideally none of these values should need to be set as uplift will always use what we believe to be sensible default values.

## commitMessage

```yaml linenums="1"
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

```yaml linenums="1"
# .uplift.yml

# Changes the commit author used by uplift when committing any staged changes.
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

## annotatedTags

```yaml linenums="1"
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
