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

func newChangelogCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Create or update a changelog with the latest semantic release",
		Long:  chlogDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Attempt to retrieve the latest 2 tags for generating a changelog entry
			tags := git.AllTags()
			if len(tags) == 1 {
				ctx.NextVersion.Raw = tags[0]
			} else if len(tags) > 1 {
				ctx.NextVersion.Raw = tags[0]
				ctx.CurrentVersion.Raw = tags[1]
			}

			if ctx.ChangelogDiff {
				// Run a condensed workflow when just calculating the diff
				return writeChangelogDiff(ctx)
			}

			return writeChangelog(ctx)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&ctx.ChangelogDiff, "diff-only", false, "output the changelog diff only")
	f.StringArrayVar(&ctx.ChangelogExcludes, "exclude", []string{}, "a list of conventional commit prefixes to exclude")

	return cmd
}

func writeChangelog(ctx *context.Context) error {
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

func writeChangelogDiff(ctx *context.Context) error {
	tsk := changelog.Task{}
	if err := skip.Running(tsk.Skip, logging.Log(tsk.String(), tsk.Run))(ctx); err != nil {
		return err
	}

	return nil
}
