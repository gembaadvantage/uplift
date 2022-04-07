# Introduction

"Semantic versioning the easy way. Powered by Conventional Commits. Built for use with CI"

Uplift is designed to simplify release management within a project. By harnessing the power of [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/), release automation can be introduced into any CI workflow. Uplift works by analysing the latest conventional commit message to identify the next semantic release. If a semantic release is identified, uplift will release the project. A release can comprise of file bumping (updating the version within a file), changelog management and tagging of a repository, all configurable through uplifts commands and/or configuration file.

Being built using Go, uplift is incredibly small and easy to install into any CI workflow. Once you are setup, you won't have to do anything again. Uplift will take care of the rest for you!
