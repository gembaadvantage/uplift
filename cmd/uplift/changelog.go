/*
Copyright (c) 2021 Gemba Advantage

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

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/middleware/logging"
	"github.com/gembaadvantage/uplift/internal/middleware/skip"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/gembaadvantage/uplift/internal/task/changelog"
	"github.com/gembaadvantage/uplift/internal/task/gitpush"
	"github.com/gembaadvantage/uplift/internal/task/lastcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextcommit"
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
	globalOptions
}

type changelogCommand struct {
	Cmd  *cobra.Command
	Opts changelogOptions
}

func newChangelogCmd(gopts globalOptions, out io.Writer) *changelogCommand {
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
			if chglogCmd.Opts.DiffOnly {
				// Run a condensed workflow when just calculating the diff
				return writeChangelogDiff(chglogCmd.Opts, out)
			}

			return writeChangelog(chglogCmd.Opts, out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&chglogCmd.Opts.DiffOnly, "diff-only", false, "output the changelog diff only")
	f.StringSliceVar(&chglogCmd.Opts.Exclude, "exclude", []string{}, "a list of conventional commit prefixes to exclude")

	chglogCmd.Cmd = cmd
	return chglogCmd
}

func writeChangelog(opts changelogOptions, out io.Writer) error {
	ctx, err := setupChangelogContext(opts, out)
	if err != nil {
		return err
	}

	tsks := []task.Runner{
		lastcommit.Task{},
		nextcommit.Task{},
		changelog.Task{},
		gitpush.Task{},
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

	tsk := changelog.Task{}
	if err := skip.Running(tsk.Skip, logging.Log(tsk.String(), tsk.Run))(ctx); err != nil {
		return err
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

	// Set all values within the context
	ctx.Debug = opts.Debug
	ctx.DryRun = opts.DryRun
	ctx.NoPush = opts.NoPush

	// Merge config and command line arguments together
	ctx.ChangelogExcludes = opts.Exclude
	ctx.ChangelogExcludes = append(ctx.ChangelogExcludes, ctx.Config.Changelog.Exclude...)

	// Attempt to retrieve the latest 2 tags for generating a changelog entry
	tags := git.AllTags()
	if len(tags) == 1 {
		ctx.NextVersion.Raw = tags[0]
	} else if len(tags) > 1 {
		ctx.NextVersion.Raw = tags[0]
		ctx.CurrentVersion.Raw = tags[1]
	}

	return ctx, nil
}
