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

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/gembaadvantage/uplift/internal/task/bump"
	"github.com/gembaadvantage/uplift/internal/task/changelog"
	"github.com/gembaadvantage/uplift/internal/task/fetchtag"
	"github.com/gembaadvantage/uplift/internal/task/gitcheck"
	"github.com/gembaadvantage/uplift/internal/task/gitcommit"
	"github.com/gembaadvantage/uplift/internal/task/gittag"
	"github.com/gembaadvantage/uplift/internal/task/gpgimport"
	"github.com/gembaadvantage/uplift/internal/task/hook/after"
	"github.com/gembaadvantage/uplift/internal/task/hook/afterbump"
	"github.com/gembaadvantage/uplift/internal/task/hook/afterchangelog"
	"github.com/gembaadvantage/uplift/internal/task/hook/aftertag"
	"github.com/gembaadvantage/uplift/internal/task/hook/before"
	"github.com/gembaadvantage/uplift/internal/task/hook/beforebump"
	"github.com/gembaadvantage/uplift/internal/task/hook/beforechangelog"
	"github.com/gembaadvantage/uplift/internal/task/hook/beforetag"
	"github.com/gembaadvantage/uplift/internal/task/nextcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextsemver"
	"github.com/gembaadvantage/uplift/internal/task/scm"
	"github.com/spf13/cobra"
)

const (
	releaseLongDesc = `Release the next semantic version of your git repository. A release consists of
a three-stage process. First, all configured files will be bumped (patched) using
the next semantic version. Second, a changelog entry containing all commits for
the latest semantic release will be created. Finally, Uplift will tag the
repository. Uplift automatically handles the staging and pushing of modified
files and the tagging of the repository with two separate git pushes. But this
behavior can be disabled to manage these actions manually.

Parts of this release process can be disabled if needed.

https://upliftci.dev/first-release/`

	releaseExamples = `
# Release the next semantic version
uplift release

# Release the next semantic version without bumping any files
uplift release --skip-bumps

# Release the next semantic version without generating a changelog
uplift release --skip-changelog

# Append a prerelease suffix to the next calculated semantic version
uplift release --prerelease beta.1

# Ensure any "v" prefix is stripped from the next calculated semantic
# version to explicitly adhere to the SemVer specification
uplift release --no-prefix`
)

type releaseOptions struct {
	FetchTags      bool
	Check          bool
	Prerelease     string
	SkipChangelog  bool
	SkipBumps      bool
	NoPrefix       bool
	Exclude        []string
	Include        []string
	Sort           string
	Multiline      bool
	SkipPrerelease bool
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
		Use:     "release",
		Short:   "Release the next semantic version of a repository",
		Long:    releaseLongDesc,
		Example: releaseExamples,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Just check if uplift would trigger a release
			if relCmd.Opts.Check {
				return checkRelease(relCmd.Opts, out)
			}

			return release(relCmd.Opts, out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&relCmd.Opts.FetchTags, "fetch-all", false, "fetch all tags from the remote repository")
	f.BoolVar(&relCmd.Opts.Check, "check", false, "check if a release will be triggered")
	f.StringVar(&relCmd.Opts.Prerelease, "prerelease", "", "append a prerelease suffix to next calculated semantic version")
	f.BoolVar(&relCmd.Opts.SkipChangelog, "skip-changelog", false, "skips the creation or amendment of a changelog")
	f.BoolVar(&relCmd.Opts.SkipBumps, "skip-bumps", false, "skips the bumping of any files")
	f.BoolVar(&relCmd.Opts.NoPrefix, "no-prefix", false, "strip the default 'v' prefix from the next calculated semantic version")
	f.StringSliceVar(&relCmd.Opts.Exclude, "exclude", []string{}, "a list of regexes for excluding conventional commits from the changelog")
	f.StringSliceVar(&relCmd.Opts.Include, "include", []string{}, "a list of regexes to cherry-pick conventional commits for the changelog")
	f.StringVar(&relCmd.Opts.Sort, "sort", "", "the sort order of commits within each changelog entry")
	f.BoolVar(&relCmd.Opts.Multiline, "multiline", false, "include multiline commit messages within changelog (skips truncation)")
	f.BoolVar(&relCmd.Opts.SkipPrerelease, "skip-changelog-prerelease", false, "skips the creation of a changelog entry for a prerelease")

	relCmd.Cmd = cmd
	return relCmd
}

func release(opts releaseOptions, out io.Writer) error {
	ctx, err := setupReleaseContext(opts, out)
	if err != nil {
		return err
	}

	tasks := []task.Runner{
		gitcheck.Task{},
		before.Task{},
		gpgimport.Task{},
		scm.Task{},
		fetchtag.Task{},
		nextsemver.Task{},
		nextcommit.Task{},
		beforebump.Task{},
		bump.Task{},
		afterbump.Task{},
		beforechangelog.Task{},
		changelog.Task{},
		afterchangelog.Task{},
		gitcommit.Task{},
		beforetag.Task{},
		gittag.Task{},
		aftertag.Task{},
		after.Task{},
	}

	return task.Execute(ctx, tasks)
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
	ctx.NoStage = opts.NoStage
	ctx.FetchTags = opts.FetchTags
	ctx.Out = out
	ctx.SkipChangelog = opts.SkipChangelog
	ctx.SkipBumps = opts.SkipBumps
	ctx.NoPrefix = opts.NoPrefix

	// Enable pre-tagging support for generating a changelog
	ctx.Changelog.PreTag = true

	// Merge config and command line arguments together
	ctx.Changelog.Multiline = opts.Multiline
	if !ctx.Changelog.Multiline && ctx.Config.Changelog != nil {
		ctx.Changelog.Multiline = ctx.Config.Changelog.Multiline
	}
	ctx.Changelog.Include = opts.Include
	ctx.Changelog.Exclude = opts.Exclude
	if ctx.Config.Changelog != nil {
		ctx.Changelog.Include = append(ctx.Changelog.Include, ctx.Config.Changelog.Include...)
		ctx.Changelog.Exclude = append(ctx.Changelog.Exclude, ctx.Config.Changelog.Exclude...)
	}
	ctx.Changelog.SkipPrerelease = opts.SkipPrerelease
	if !ctx.Changelog.SkipPrerelease && ctx.Config.Changelog != nil {
		ctx.Changelog.SkipPrerelease = ctx.Config.Changelog.SkipPrerelease
	}

	// By default ensure the ci(uplift): commits are excluded also
	ctx.Changelog.Exclude = append(ctx.Changelog.Exclude, "ci(uplift):")

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

func checkRelease(opts releaseOptions, out io.Writer) error {
	ctx, err := setupReleaseContext(opts, out)
	if err != nil {
		return err
	}

	tasks := []task.Runner{
		nextsemver.Task{},
	}

	if err := task.Execute(ctx, tasks); err != nil {
		return err
	}

	if ctx.NoVersionChanged {
		return errors.New("no release detected")
	}

	return nil
}
