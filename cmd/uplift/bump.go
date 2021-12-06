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
	"github.com/gembaadvantage/uplift/internal/middleware/logging"
	"github.com/gembaadvantage/uplift/internal/middleware/skip"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/gembaadvantage/uplift/internal/task/bump"
	"github.com/gembaadvantage/uplift/internal/task/currentversion"
	"github.com/gembaadvantage/uplift/internal/task/gitpush"
	"github.com/gembaadvantage/uplift/internal/task/lastcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextversion"
	"github.com/spf13/cobra"
)

const (
	bumpDesc = `Bumps the semantic version within files in your git repository. The
version bump is based on the conventional commit message from the last commit.
Uplift can bump the version in any file using regex pattern matching`
)

func newBumpCmd(ctx *context.Context) *cobra.Command {
	var pre string

	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Bump the semantic version within files",
		Long:  bumpDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle prerelease suffix if one is provided
			if pre != "" {
				var err error
				if ctx.Prerelease, ctx.Metadata, err = semver.ParsePrerelease(pre); err != nil {
					return err
				}
			}

			return bumpFiles(ctx)
		},
	}

	f := cmd.Flags()
	f.StringVar(&pre, "prerelease", "", "append a prerelease suffix to next calculated semantic version")

	return cmd
}

func bumpFiles(ctx *context.Context) error {
	tsks := []task.Runner{
		lastcommit.Task{},
		currentversion.Task{},
		nextversion.Task{},
		nextcommit.Task{},
		bump.Task{},
		gitpush.Task{},
	}

	for _, tsk := range tsks {
		if err := skip.Running(tsk.Skip, logging.Log(tsk.String(), tsk.Run))(ctx); err != nil {
			return err
		}
	}

	return nil
}
