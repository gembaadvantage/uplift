/*
Copyright (c) 2022 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/middleware/logging"
	"github.com/gembaadvantage/uplift/internal/middleware/skip"
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
	"github.com/spf13/cobra"
)

const (
	chlogDesc = `Create or update an existing changelog with an entry for
the latest semantic release. For a first release, all commits
between the latest tag and trunk will be written to the
changelog. Subsequent entries will contain only commits between 
release tags`
)

type changelogOptions struct {
	DiffOnly bool
	Exclude  []string
	Include  []string
	All      bool
	Sort     string
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
		Use:   "changelog",
		Short: "Create or update a changelog with the latest semantic release",
		Long:  chlogDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	chglogCmd.Cmd = cmd
	return chglogCmd
}

func writeChangelog(opts changelogOptions, out io.Writer) error {
	ctx, err := setupChangelogContext(opts, out)
	if err != nil {
		return err
	}

	tsks := []task.Runner{
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

	for _, tsk := range tsks {
		if err := skip.Running(tsk.Skip, logging.Log(tsk.String(), tsk.Run))(ctx); err != nil {
			return err
		}
	}

	return nil
}

func writeChangelogDiff(opts changelogOptions, out io.Writer) error {
	ctx, err := setupChangelogContext(opts, out)
	if err != nil {
		return err
	}

	tsks := []task.Runner{
		gitcheck.Task{},
		before.Task{},
		scm.Task{},
		changelog.Task{},
		after.Task{},
	}

	for _, tsk := range tsks {
		if err := skip.Running(tsk.Skip, logging.Log(tsk.String(), tsk.Run))(ctx); err != nil {
			return err
		}
	}

	return nil
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

	// Sort order provided as a command-line flag takes precedence
	ctx.Changelog.Sort = opts.Sort
	if ctx.Changelog.Sort == "" {
		ctx.Changelog.Sort = strings.ToLower(cfg.Changelog.Sort)
	}

	// Merge config and command line arguments together
	ctx.Changelog.Include = append(opts.Include, ctx.Config.Changelog.Include...)
	ctx.Changelog.Exclude = append(opts.Exclude, ctx.Config.Changelog.Exclude...)

	// By default ensure the ci(uplift): commits are excluded also
	ctx.Changelog.Exclude = append(ctx.Changelog.Exclude, `ci\(uplift\)`)

	if !ctx.Changelog.All {
		// Attempt to retrieve the latest 2 tags for generating a changelog entry
		tags := git.AllTags()
		if len(tags) == 1 {
			ctx.NextVersion.Raw = tags[0].Ref
		} else if len(tags) > 1 {
			ctx.NextVersion.Raw = tags[0].Ref
			ctx.CurrentVersion.Raw = tags[1].Ref
		}
	}

	// Handle git config. Command line flag takes precedences
	ctx.IgnoreDetached = opts.IgnoreDetached
	if !ctx.IgnoreDetached {
		ctx.IgnoreDetached = ctx.Config.Git.IgnoreDetached
	}

	ctx.IgnoreShallow = opts.IgnoreShallow
	if !ctx.IgnoreShallow {
		ctx.IgnoreShallow = ctx.Config.Git.IgnoreShallow
	}

	return ctx, nil
}
