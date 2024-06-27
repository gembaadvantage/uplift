package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/gembaadvantage/uplift/internal/task/changelog"
	"github.com/gembaadvantage/uplift/internal/task/gitcheck"
	"github.com/gembaadvantage/uplift/internal/task/gitcommit"
	"github.com/gembaadvantage/uplift/internal/task/hook/after"
	"github.com/gembaadvantage/uplift/internal/task/hook/afterchangelog"
	"github.com/gembaadvantage/uplift/internal/task/hook/before"
	"github.com/gembaadvantage/uplift/internal/task/hook/beforechangelog"
	"github.com/gembaadvantage/uplift/internal/task/nextcommit"
	"github.com/gembaadvantage/uplift/internal/task/scm"
	git "github.com/purpleclay/gitz"
	"github.com/spf13/cobra"
)

const (
	changelogLongDesc = `Scans the git log for the latest semantic release and generates a changelog
entry. If this is a first release, all commits between the last release (or
identifiable tag) and the repository trunk will be written to the changelog.
Any subsequent entry within the changelog will only contain commits between
the latest set of tags. Basic customization is supported. Optionally commits
can be explicitly included or excluded from the entry and sorted in ascending
or descending order. Uplift automatically handles the staging and pushing of
changes to the CHANGELOG.md file to the git remote, but this behavior can be
disabled, to manage this action manually.

Uplift bases its changelog format on the Keep a Changelog specification:

https://keepachangelog.com/en/1.0.0/`

	changelogExamples = `
# Generate the next changelog entry for the latest semantic release
uplift changelog

# Generate a changelog for the entire history of the repository
uplift changelog --all

# Generate the next changelog entry and write it to stdout
uplift changelog --diff-only

# Generate the next changelog entry by exclude any conventional commits
# with the ci, chore or test prefixes
uplift changelog --exclude "^ci,^chore,^test"

# Generate the next changelog entry with commits that only include the
# following scope
uplift changelog --include "^.*\(scope\)"

# Generate the next changelog entry but do not stage or push any changes
# back to the git remote
uplift changelog --no-stage

# Generate a changelog with multiline commit messages
uplift changelog --multiline

# Generate a changelog trimming any lines preceding the conventional commit type
uplift changelog --trim-header

# Generate a changelog with prerelease tags being skipped
uplift changelog --skip-prerelease`
)

type changelogOptions struct {
	DiffOnly       bool
	Exclude        []string
	Include        []string
	All            bool
	Sort           string
	Multiline      bool
	SkipPrerelease bool
	TrimHeader     bool
	*globalOptions
}

type changelogCommand struct {
	Cmd  *cobra.Command
	Opts changelogOptions
}

func newChangelogCmd(gopts *globalOptions, out io.Writer) *changelogCommand {
	chglogCmd := &changelogCommand{
		Opts: changelogOptions{
			globalOptions: gopts,
		},
	}

	cmd := &cobra.Command{
		Use:     "changelog",
		Short:   "Create or update a changelog with the latest semantic release",
		Long:    changelogLongDesc,
		Example: changelogExamples,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			// Always lowercase sort
			chglogCmd.Opts.Sort = strings.ToLower(chglogCmd.Opts.Sort)

			if chglogCmd.Opts.DiffOnly {
				// Run a condensed workflow when just calculating the diff
				return writeChangelogDiff(chglogCmd.Opts, out)
			}

			return writeChangelog(chglogCmd.Opts, out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&chglogCmd.Opts.DiffOnly, "diff-only", false, "output the changelog diff only")
	f.BoolVar(&chglogCmd.Opts.All, "all", false, "generate a changelog from the entire history of this repository")
	f.StringSliceVar(&chglogCmd.Opts.Exclude, "exclude", []string{}, "a list of regexes for excluding conventional commits from the changelog")
	f.StringSliceVar(&chglogCmd.Opts.Include, "include", []string{}, "a list of regexes to cherry-pick conventional commits for the changelog")
	f.StringVar(&chglogCmd.Opts.Sort, "sort", "", "the sort order of commits within each changelog entry")
	f.BoolVar(&chglogCmd.Opts.Multiline, "multiline", false, "include multiline commit messages within changelog (skips truncation)")
	f.BoolVar(&chglogCmd.Opts.SkipPrerelease, "skip-prerelease", false, "skips the creation of a changelog entry for a prerelease")
	f.BoolVar(&chglogCmd.Opts.TrimHeader, "trim-header", false, "strip any lines preceding the conventional commit type in the commit message")

	chglogCmd.Cmd = cmd
	return chglogCmd
}

func writeChangelog(opts changelogOptions, out io.Writer) error {
	ctx, err := setupChangelogContext(opts, out)
	if err != nil {
		return err
	}

	tasks := []task.Runner{
		gitcheck.Task{},
		before.Task{},
		scm.Task{},
		nextcommit.Task{},
		beforechangelog.Task{},
		changelog.Task{},
		afterchangelog.Task{},
		gitcommit.Task{},
		after.Task{},
	}

	return task.Execute(ctx, tasks)
}

func writeChangelogDiff(opts changelogOptions, out io.Writer) error {
	ctx, err := setupChangelogContext(opts, out)
	if err != nil {
		return err
	}

	tasks := []task.Runner{
		gitcheck.Task{},
		before.Task{},
		scm.Task{},
		changelog.Task{},
		after.Task{},
	}

	return task.Execute(ctx, tasks)
}

func setupChangelogContext(opts changelogOptions, out io.Writer) (*context.Context, error) {
	cfg, err := loadConfig(opts.ConfigDir)
	if err != nil {
		fmt.Printf("failed to load uplift config. %v", err)
		return nil, err
	}
	ctx := context.New(cfg, out)
	ctx.Out = out

	// Set all values within the context
	ctx.Debug = opts.Debug
	ctx.DryRun = opts.DryRun
	ctx.NoPush = opts.NoPush
	ctx.NoStage = opts.NoStage
	ctx.Changelog.DiffOnly = opts.DiffOnly
	ctx.Changelog.All = opts.All
	ctx.Changelog.Multiline = opts.Multiline
	if !ctx.Changelog.Multiline && ctx.Config.Changelog != nil {
		ctx.Changelog.Multiline = ctx.Config.Changelog.Multiline
	}
	ctx.Changelog.SkipPrerelease = opts.SkipPrerelease
	if !ctx.Changelog.SkipPrerelease && ctx.Config.Changelog != nil {
		ctx.Changelog.SkipPrerelease = ctx.Config.Changelog.SkipPrerelease
	}

	ctx.Changelog.TrimHeader = opts.TrimHeader
	if !ctx.Changelog.TrimHeader && ctx.Config.Changelog != nil {
		ctx.Changelog.TrimHeader = ctx.Config.Changelog.TrimHeader
	}

	// Sort order provided as a command-line flag takes precedence
	ctx.Changelog.Sort = opts.Sort
	if ctx.Changelog.Sort == "" && cfg.Changelog != nil {
		ctx.Changelog.Sort = strings.ToLower(cfg.Changelog.Sort)
	}

	// Merge config and command line arguments together
	ctx.Changelog.Include = opts.Include
	ctx.Changelog.Exclude = opts.Exclude
	if ctx.Config.Changelog != nil {
		ctx.Changelog.Include = append(ctx.Changelog.Include, ctx.Config.Changelog.Include...)
		ctx.Changelog.Exclude = append(ctx.Changelog.Exclude, ctx.Config.Changelog.Exclude...)
	}

	// By default ensure the ci(uplift): commits are excluded also
	ctx.Changelog.Exclude = append(ctx.Changelog.Exclude, `ci\(uplift\)`)

	if !ctx.Changelog.All {
		// Attempt to retrieve the latest 2 tags for generating a changelog entry
		tags, err := ctx.GitClient.Tags(git.WithShellGlob("*.*.*"),
			git.WithSortBy(git.CreatorDateDesc, git.VersionDesc))
		if err != nil {
			return nil, err
		}

		if len(tags) == 1 {
			ctx.NextVersion.Raw = tags[0]
		} else if len(tags) > 1 {
			ctx.NextVersion.Raw = tags[0]
			ctx.CurrentVersion.Raw = tags[1]
		}
	}

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
