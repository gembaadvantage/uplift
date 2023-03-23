# Configure GitLab pipelines with Push Options

Since Git version `2.10`, a client could send arbitrary strings to a server (_remote_) using [push options](https://git-scm.com/docs/git-push#Documentation/git-push.txt--oltoptiongt) (`--push-option`). GitLab has utilized [push options](https://docs.gitlab.com/ee/user/project/push_options.html) to configure pipeline behavior since version `11.7`. This type of configuration opens up many possibilities for configuring your CI/CD workflow through uplift.

## Skip Rebuilding Main Branch

During a release, uplift may commit and push some staged files (_bumped files and changelog_) before tagging the project with the next semantic version. GitLab's default behavior is to spawn two pipelines, one for the main branch and another for the tag. The former may be deemed surplus to requirements and can be skipped using the `ci.skip` push option:

```{ .yaml .annotate linenums="1" }
git:
  pushOptions:
    - option: ci.skip
      skipTag: true
```

## Dynamically update CI/CD Variables

Using environment variables to configure pipelines is a common practice within GitLab. It may be desirable in certain conditions to dynamically change their values, and GitLab provides a `ci.variable` push option for this exact purpose.

```{ .yaml .annotate linenums="1" }
git:
  pushOptions:
    - option: ci.variable="CODE_QUALITY_DISABLED=true"
      skipBranch: true
```
