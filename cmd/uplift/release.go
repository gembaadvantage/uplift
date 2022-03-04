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
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/middleware/logging"
	"github.com/gembaadvantage/uplift/internal/middleware/skip"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/gembaadvantage/uplift/internal/task/bump"
	"github.com/gembaadvantage/uplift/internal/task/changelog"
	"github.com/gembaadvantage/uplift/internal/task/currentversion"
	"github.com/gembaadvantage/uplift/internal/task/fetchtag"
	"github.com/gembaadvantage/uplift/internal/task/gitcommit"
	"github.com/gembaadvantage/uplift/internal/task/gittag"
	"github.com/gembaadvantage/uplift/internal/task/lastcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextversion"
	"github.com/gembaadvantage/uplift/internal/task/scm"
	"github.com/spf13/cobra"
)

const (
	releaseDesc = `Release the next semantic version of your git repository. A release
will automatically bump any files and tag the associated commit with 
the required semantic version`
)

type releaseOptions struct {
	FetchTags  bool
	Check      bool
	Prerelease string
	*globalOptions
}

type releaseCommand struct {
	Cmd  *cobra.Command
	Opts releaseOptions
}

func newReleaseCmd(gopts *globalOptions, out io.Writer) *releaseCommand {
	relCmd := &releaseCommand{
		Opts: releaseOptions{
			globalOptions: gopts,
		},
	}

	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release the next semantic version of a repository",
		Long:  releaseDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Just check if uplift would trigger a release
			if relCmd.Opts.Check {
				return checkRelease()
			}

			return release(relCmd.Opts, out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&relCmd.Opts.FetchTags, "fetch-all", false, "fetch all tags from the remote repository")
	f.BoolVar(&relCmd.Opts.Check, "check", false, "check if a release will be triggered")
	f.StringVar(&relCmd.Opts.Prerelease, "prerelease", "", "append a prerelease suffix to next calculated semantic version")

	relCmd.Cmd = cmd
	return relCmd
}

func release(opts releaseOptions, out io.Writer) error {
	ctx, err := setupReleaseContext(opts, out)
	if err != nil {
		return err
	}

	tsks := []task.Runner{
		fetchtag.Task{},
		lastcommit.Task{},
		currentversion.Task{},
		nextversion.Task{},
		nextcommit.Task{},
		scm.Task{},
		bump.Task{},
		changelog.Task{},
		gitcommit.Task{},
		gittag.Task{},
	}

	for _, tsk := range tsks {
		if err := skip.Running(tsk.Skip, logging.Log(tsk.String(), tsk.Run))(ctx); err != nil {
			return err
		}
	}

	return nil
}

func setupReleaseContext(opts releaseOptions, out io.Writer) (*context.Context, error) {
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
	ctx.Out = out

	// Enable pre-tagging support for generating a changelog
	ctx.Changelog.PreTag = true

	// Handle prerelease suffix if one is provided
	if opts.Prerelease != "" {
		var err error
		if ctx.Prerelease, ctx.Metadata, err = semver.ParsePrerelease(opts.Prerelease); err != nil {
			return nil, err
		}
	}

	return ctx, nil
}

func checkRelease() error {
	return logging.Log("check release", func(ctx *context.Context) error {
		cd, err := git.LatestCommit()
		if err != nil {
			return err
		}

		log.WithField("message", cd.Message).Info("retrieved latest commit")

		inc := semver.ParseCommit(cd.Message)
		if inc == semver.NoIncrement {
			log.Info("nothing to release")
			return errors.New("no release would be triggered for this commit")
		}

		log.WithField("increment", strings.ToLower(string(inc))).Info("detected releasable commit")
		return nil
	})(&context.Context{})
}
