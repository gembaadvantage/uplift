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

package context

import (
	ctx "context"
	"io"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/semver"
	git "github.com/purpleclay/gitz"
)

// Context provides a way to share common state across tasks
type Context struct {
	ctx.Context
	Changelog                Changelog
	CommitDetails            git.CommitDetails
	Config                   config.Uplift
	CurrentVersion           semver.Version
	IncludeArtifacts         []string
	DryRun                   bool
	Debug                    bool
	FetchTags                bool
	FilterOnPrerelease       bool
	GitClient                *git.Client
	IgnoreDetached           bool
	IgnoreExistingPrerelease bool
	IgnoreShallow            bool
	Prerelease               string
	Metadata                 string
	NextVersion              semver.Version
	NoPrefix                 bool
	NoVersionChanged         bool
	NoPush                   bool
	NoStage                  bool
	Out                      io.Writer
	PrintCurrentTag          bool
	PrintNextTag             bool
	SCM                      SCM
	SkipBumps                bool
	SkipChangelog            bool
}

// SCMProvider is used for identifying the source code management tool used
// by the current git repository
type SCMProvider string

const (
	GitHub       SCMProvider = "GitHub"
	GitLab       SCMProvider = "GitLab"
	CodeCommit   SCMProvider = "CodeCommit"
	Gitea        SCMProvider = "Gitea"
	Unrecognised SCMProvider = "Unrecognised"
)

// SCM provides details about the SCM provider of a repository
type SCM struct {
	Provider  SCMProvider
	TagURL    string
	CommitURL string
}

// Changelog provides details about how the changelog should be managed
// for the current repository
type Changelog struct {
	All            bool
	DiffOnly       bool
	Exclude        []string
	Include        []string
	Sort           string
	PreTag         bool
	Multiline      bool
	SkipPrerelease bool
	TrimHeader     bool
}

// New constructs a context that captures both runtime configuration and
// user defined runtime options
func New(cfg config.Uplift, out io.Writer) *Context {
	// TODO: this should throw an error if required
	gc, _ := git.NewClient()

	return &Context{
		Context:   ctx.Background(),
		Config:    cfg,
		GitClient: gc,
		Out:       out,
		SCM: SCM{
			Provider: Unrecognised,
		},
		IncludeArtifacts: IncludeArtifacts(cfg),
	}
}

// For nil safe object getting
func IncludeArtifacts(c config.Uplift) []string {
	if c.Git == nil {
		return nil
	}

	return c.Git.IncludeArtifacts
}
