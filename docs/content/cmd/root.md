---
sidebar_position: 1
---

# uplift

Semantic versioning the easy way. Powered by Conventional Commits. Built for use with CI.

```sh
uplift [COMMAND]
```

## Global Flags

### --config-dir

Provide a custom path to a directory containing your uplift configuration file. By default uplift will look in the current directory where it was run.

### --debug

Turn on extra debug output. Good for diving into the details of how uplift works and for reporting any issues that you discover.

### --dry-run

Run uplift without making any changes. A good way for exploring how uplift works.

### --no-push

Prevents uplift from pushing any changes back to your git remote. Any changes made by uplift will remain locally staged.

### --silent

Peace and quiet! Stops uplift from logging anything. A great option when combining uplift with any custom shell scripts.

### --ignore-detached

Suppress the git detached HEAD check within uplift. I have it all under control.

### --ignore-shallow

Suppress the git shallow clone check within uplift. I have it all under control.
