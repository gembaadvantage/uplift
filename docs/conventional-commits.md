# Why use Conventional Commits?

[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) is a specification designed to introduce human and machine-readable meaning to commit messages, enabling automated tooling such as Uplift for managing releases. A user prefixes their commit with a type to describe their intent, and these labels form a direct relationship with [Semantic Versioning](https://semver.org/). The specification adopts the Angular convention, and so does Uplift.

## Semantic Versioning Types

- `fix:` A bug fix triggers a patch semantic version bump `0.1.0 ~> 0.1.1`
- `feat:` A new feature triggers a minor semantic version bump `0.1.0 ~> 0.2.0`
- `feat{==!==}:` A breaking change triggers a major semantic version bump `0.1.0 ~> 1.0.0`[^1]

## Additional Angular Types

Uplift supports all the additional Angular types, `chore:`, `ci:`, `docs:`, `style:`, `refactor:`, `perf:` and `test:`.

## How Uplift Scans Commits

When determining the next semantic version, all commit messages for a release are scanned for the highest possible increment (Patch, Minor or Major).

```text
docs: add documentation for new exciting feature
ci: build and test documentation within pipeline
feat: shiny new feature                           <-- largest increment
fix: fixed another bug found by user
fix: fixed bug found by user
```

In the above example, if the latest tag were `0.1.0` it would be incremented to `0.2.0`.

[^1]: Users can also add a `BREAKING CHANGE` footer to their commit message.
