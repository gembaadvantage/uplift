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
	"github.com/gembaadvantage/uplift/internal/middleware/logging"
	"github.com/gembaadvantage/uplift/internal/middleware/skip"
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
	"github.com/spf13/cobra"
)

const (
	tagDesc = `Tags a git repository with the next semantic version. The tag
is based on the conventional commit message from the last commit.`
)

var (
	tagRepoPipeline = []task.Runner{
		before.Task{},
		gitcheck.Task{},
		fetchtag.Task{},
		nextsemver.Task{},
		nextcommit.Task{},
		beforetag.Task{},
		gittag.Task{},
		aftertag.Task{},
		after.Task{},
	}

	printTagPipeline = []task.Runner{
		before.Task{},
		gitcheck.Task{},
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
		Use:   "tag",
		Short: "Tag a git repository with the next semantic version",
		Long:  tagDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	tsks := tagRepoPipeline

	// Switch pipeline if either the current or next tag is to be printed only
	if ctx.PrintCurrentTag || ctx.PrintNextTag {
		tsks = printTagPipeline
	}

	for _, tsk := range tsks {
		if err := skip.Running(tsk.Skip, logging.Log(tsk.String(), tsk.Run))(ctx); err != nil {
			return err
		}
	}

	return nil
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
