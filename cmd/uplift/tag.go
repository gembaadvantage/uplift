package main

import (
	"fmt"
	"io"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/gembaadvantage/uplift/internal/task/fetchtag"
	"github.com/gembaadvantage/uplift/internal/task/gitcheck"
	"github.com/gembaadvantage/uplift/internal/task/gittag"
	"github.com/gembaadvantage/uplift/internal/task/hook/after"
	"github.com/gembaadvantage/uplift/internal/task/hook/aftertag"
	"github.com/gembaadvantage/uplift/internal/task/hook/before"
	"github.com/gembaadvantage/uplift/internal/task/hook/beforetag"
	"github.com/gembaadvantage/uplift/internal/task/nextcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextsemver"
	git "github.com/purpleclay/gitz"
	"github.com/spf13/cobra"
)

const (
	tagLongDesc = `Generates a new git tag by scanning the git log of a repository for any
conventional commits since the last release (or identifiable tag). When
examining the git log, Uplift will always calculate the next semantic version
based on the most significant detected increment. Uplift automatically handles
the creation and pushing of a new git tag to the remote, but this behavior can
be disabled, to manage this action manually.

Conventional Commits is a set of rules for creating an explicit commit history,
which makes building automation tools like Uplift much easier. Uplift adheres to
v1.0.0 of the specification:

https://www.conventionalcommits.org/en/v1.0.0`

	tagExamples = `
# Tag the repository with the next calculated semantic version
uplift tag

# Identify the next semantic version and write to stdout. Repository is
# not tagged
uplift tag --next --silent

# Identify the current semantic version and write to stdout. Repository
# is not tagged
uplift tag --current

# Identify the current and next semantic versions and write both to stdout.
# Repository is not tagged
uplift tag --current --next --silent

# Ensure the calculated version explicitly aheres to the SemVer specification
# by stripping the "v" prefix from the generated tag
uplift tag --no-prefix

# Append a prerelease suffix to the next calculated semantic version
uplift tag --prerelease beta.1

# Tag the repository with the next calculated semantic version, but do not
# push the tag to the remote
uplift tag --no-push`
)

var (
	tagRepoPipeline = []task.Runner{
		gitcheck.Task{},
		before.Task{},
		fetchtag.Task{},
		nextsemver.Task{},
		nextcommit.Task{},
		beforetag.Task{},
		gittag.Task{},
		aftertag.Task{},
		after.Task{},
	}

	printNextTagPipeline = []task.Runner{
		gitcheck.Task{},
		before.Task{},
		fetchtag.Task{},
		nextsemver.Task{},
		beforetag.Task{},
		gittag.Task{},
		aftertag.Task{},
		after.Task{},
	}
)

type tagOptions struct {
	FetchTags       bool
	PrintCurrentTag bool
	PrintNextTag    bool
	Prerelease      string
	NoPrefix        bool
	*globalOptions
}

type tagCommand struct {
	Cmd  *cobra.Command
	Opts tagOptions
}

func newTagCmd(gopts *globalOptions, out io.Writer) *tagCommand {
	tagCmd := &tagCommand{
		Opts: tagOptions{
			globalOptions: gopts,
		},
	}

	cmd := &cobra.Command{
		Use:     "tag",
		Short:   "Tag a git repository with the next calculated semantic version",
		Long:    tagLongDesc,
		Example: tagExamples,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// If only the current tag is to be printed, skip running a pipeline
			// and just retrieve and print the latest tag
			if tagCmd.Opts.PrintCurrentTag && !tagCmd.Opts.PrintNextTag {
				gc, err := git.NewClient()
				if err != nil {
					return err
				}

				tags, _ := gc.Tags(git.WithShellGlob("*.*.*"),
					git.WithSortBy(git.CreatorDateDesc, git.VersionDesc),
					git.WithCount(1))
				if len(tags) == 1 {
					fmt.Fprint(out, tags[0])
				}
				return nil
			}

			return tagRepo(tagCmd.Opts, out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&tagCmd.Opts.PrintCurrentTag, "current", false, "output the current tag")
	f.BoolVar(&tagCmd.Opts.FetchTags, "fetch-all", false, "fetch all tags from the remote repository")
	f.BoolVar(&tagCmd.Opts.PrintNextTag, "next", false, "output the next tag")
	f.BoolVar(&tagCmd.Opts.NoPrefix, "no-prefix", false, "strip the default 'v' prefix from the next calculated semantic version")
	f.StringVar(&tagCmd.Opts.Prerelease, "prerelease", "", "append a prerelease suffix to next calculated semantic version")

	tagCmd.Cmd = cmd
	return tagCmd
}

func tagRepo(opts tagOptions, out io.Writer) error {
	ctx, err := setupTagContext(opts, out)
	if err != nil {
		return err
	}

	tasks := tagRepoPipeline

	if ctx.PrintNextTag {
		tasks = printNextTagPipeline
	}

	return task.Execute(ctx, tasks)
}

func setupTagContext(opts tagOptions, out io.Writer) (*context.Context, error) {
	cfg, err := loadConfig(opts.ConfigDir)
	if err != nil {
		fmt.Printf("failed to load uplift config. %v", err)
		return nil, err
	}
	ctx := context.New(cfg, out)

	// Set all values within the context
	ctx.Debug = opts.Debug
	ctx.DryRun = opts.DryRun
	ctx.NoPush = opts.NoPush
	ctx.FetchTags = opts.FetchTags
	ctx.PrintCurrentTag = opts.PrintCurrentTag
	ctx.PrintNextTag = opts.PrintNextTag
	ctx.Out = out
	ctx.NoPrefix = opts.NoPrefix

	// Handle prerelease suffix if one is provided
	if opts.Prerelease != "" {
		var err error
		if ctx.Prerelease, ctx.Metadata, err = semver.ParsePrerelease(opts.Prerelease); err != nil {
			return nil, err
		}
	}
	ctx.IgnoreExistingPrerelease = opts.IgnoreExistingPrerelease
	ctx.FilterOnPrerelease = opts.FilterOnPrerelease

	// Handle git config. Command line flag takes precedences
	ctx.IgnoreDetached = opts.IgnoreDetached
	if !ctx.IgnoreDetached && ctx.Config.Git != nil {
		ctx.IgnoreDetached = ctx.Config.Git.IgnoreDetached
	}

	ctx.IgnoreShallow = opts.IgnoreShallow
	if !ctx.IgnoreShallow && ctx.Config.Git != nil {
		ctx.IgnoreShallow = ctx.Config.Git.IgnoreShallow
	}

	return ctx, nil
}
