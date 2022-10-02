# Changing the Commit Details

If you don't want to use the default commit from the Uplift-Bot, you are free to replace both the author and commit message with anything you like:

```yaml linenums="1"
commitAuthor:
  name: "joe.bloggs"
  email: "joe.bloggs@gmail.com"

commitMessage: "chore(release): a custom release message"
```

!!!tip "Uplift uses your GPG Identity"

    If you have imported a GPG key, Uplift will always use the keys user identity over any other configuration
