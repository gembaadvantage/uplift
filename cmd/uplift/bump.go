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

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/gembaadvantage/uplift/internal/task/bump"
	"github.com/gembaadvantage/uplift/internal/task/gitcheck"
	"github.com/gembaadvantage/uplift/internal/task/gitcommit"
	"github.com/gembaadvantage/uplift/internal/task/gpgimport"
	"github.com/gembaadvantage/uplift/internal/task/hook/after"
	"github.com/gembaadvantage/uplift/internal/task/hook/afterbump"
	"github.com/gembaadvantage/uplift/internal/task/hook/before"
	"github.com/gembaadvantage/uplift/internal/task/hook/beforebump"
	"github.com/gembaadvantage/uplift/internal/task/nextcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextsemver"
	"github.com/spf13/cobra"
)

const (
	bumpLongDesc = `Calculates the next semantic version based on the conventional commits since the
last release (or identifiable tag) and bumps (or patches) a configurable set of
files with said version. JSON Path or Regex Pattern matching is supported when
scanning files for an existing semantic version. Uplift automatically handles
the staging and pushing of modified files to the git remote, but this behavior
can be disabled, to manage this action manually.

Configuring a bump requires an Uplift configuration file to exist within the
root of your project:

https://upliftci.dev/bumping-files/`

	bumpExamples = `
# Bump (patch) all configured files with the next calculated semantic version
uplift bump

# Append a prerelease suffix to the next calculated semantic version
uplift bump --prerelease beta.1

# Bump (patch) all configured files but do not stage or push any changes
# back to the git remote
uplift bump --no-stage`
)

type bumpOptions struct {
	Prerelease string
	*globalOptions
}

type bumpCommand struct {
	Cmd  *cobra.Command
	Opts bumpOptions
}

func newBumpCmd(gopts *globalOptions, out io.Writer) *bumpCommand {
	bmpCmd := &bumpCommand{
		Opts: bumpOptions{
			globalOptions: gopts,
		},
	}

	cmd := &cobra.Command{
		Use:     "bump",
		Short:   "Bump the semantic version within files",
		Long:    bumpLongDesc,
		Example: bumpExamples,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return bumpFiles(bmpCmd.Opts, out)
		},
	}

	f := cmd.Flags()
	f.StringVar(&bmpCmd.Opts.Prerelease, "prerelease", "", "append a prerelease suffix to next calculated semantic version")

	bmpCmd.Cmd = cmd
	return bmpCmd
}

func bumpFiles(opts bumpOptions, out io.Writer) error {
	ctx, err := setupBumpContext(opts, out)
	if err != nil {
		return err
	}

	tasks := []task.Runner{
		gitcheck.Task{},
		before.Task{},
		gpgimport.Task{},
		nextsemver.Task{},
		nextcommit.Task{},
		beforebump.Task{},
		bump.Task{},
		afterbump.Task{},
		gitcommit.Task{},
		after.Task{},
	}

	return task.Execute(ctx, tasks)
}

func setupBumpContext(opts bumpOptions, out io.Writer) (*context.Context, error) {
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
	ctx.NoStage = opts.NoStage
	ctx.Out = out

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
