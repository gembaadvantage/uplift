---
sidebar_position: 2
---

# uplift tag

Tags a git repository with the next semantic version. The next tag is calculated using the conventional commit message from the last commit.

```sh
uplift tag [FLAGS]
```

## Flags

### --fetch-all

Ensure all tags associated with the git repository are fetched before identifying the next tag.

```sh
uplift tag --fetch-all
```

### --next

Identify and output the next tag based on the conventional commit message of the last commit. Automatically disables tagging of the git repository. Useful when combining uplift with other external tools.

```sh
$ uplift tag --next
   • latest commit
      • retrieved latest commit   author=joe.bloggs email=joe.bloggs@example.com message=feat: this is a new feature
   • current version
      • identified version        current=0.1.2
   • next version
      • identified next version   commit=feat: this is a new feature current=0.1.2 next=0.2.0
   • next commit
      • changes will be committed with email=joe.bloggs@example.com message=ci(uplift): uplifted for version 0.2.0 name=joe.bloggs
   • git tag
      • identified next tag       tag=0.2.0
0.2.0
```

Easily capture the tag within an environment variable for use within a custom script:

```sh
NEXT_TAG=$(uplift tag --next --silent)
```

### --no-prefix

Strips the default `v` prefix from the next calculated semantic version.

```sh
uplift release --no-prefix
```

!!!tip "Should only need to use this once"

    Once a repository has been tagged with either scheme, e.g. `1.0.0` or `v1.0.0`, uplift will continue using it. This flag is most useful when tagging your repository for the first time.

### --prerelease

Append a prerelease suffix to the next calculated semantic version and use it as a prerelease tag. Supporting the [Semver 2.0.0](https://semver.org/) specification, additional labels can be provided for both the prerelease and metadata parts:

- 1.0.0`-beta.1`
- 1.1.0`-beta.1+20220312`

```sh
$ uplift tag --prerelease beta.1
   • latest commit
      • retrieved latest commit   author=joe.bloggs email=joe.bloggs@example.com message=feat: this is a new feature
   • current version
      • identified version        current=0.1.2
   • next version
      • identified next version   commit=feat: this is a new feature current=0.1.2 next=0.2.0-beta.1
   • next commit
      • changes will be committed with email=joe.bloggs@example.com message=ci(uplift): uplifted for version 0.2.0-beta.1 name=joe.bloggs
   • git tag
      • identified next tag       tag=0.2.0-beta.1
      • tagged repository with standard tag
```
