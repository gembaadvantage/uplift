# Configure GitLab pipelines with Push Options

Since Git version 2.10, a client could send arbitrary strings to a server (_remote_) using [push options](https://git-scm.com/docs/git-push#Documentation/git-push.txt--oltoptiongt) (`--push-option`). GitLab has utilized [push options](https://docs.gitlab.com/ee/user/project/push_options.html) to configure pipeline behavior since version 11.7. This type of configuration opens up many possibilities for configuring your CI/CD workflow through uplift.

## Skip Rebuilding Main Branch

During a release, uplift may commit and push some staged files (_bumped files and changelog_) before tagging the project with the next semantic version. This can result in two pipelines being generated within GitLab. Rebuilding the entire main branch again may be deemed surplus to requirements, with just the tag pipeline needed. In this instance, the `ci.skip` push option can be applied to the branch push only.

```{ .yaml .annotate linenums="1" }
git:
  pushOptions:
    - option: ci.skip
      skipTag: true
```

## Dynamically update CI/CD Variables

GitLab provides AutoDevOps jobs that run across all pipelines types, an example being the Code Quality job. After uplift has tagged a new release, it may be deemed surplus to requirements to rerun the code quality job. The `CODE_QUALITY_DISABLED` CI variable can be dynamically set on the tag pipeline.

```{ .yaml .annotate linenums="1" }
git:
  pushOptions:
    - option: ci.variable="CODE_QUALITY_DISABLED=true"
      skipBranch: true
```
