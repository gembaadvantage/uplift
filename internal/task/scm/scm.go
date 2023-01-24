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

package scm

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/codecommit-sign/pkg/translate"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
)

// Task that determines the SCM provider of a repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "detecting scm provider"
}

// Skip is disabled for this task
func (t Task) Skip(ctx *context.Context) bool {
	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	log.Debug("parsing repository remote")
	rem, err := git.Remote()
	if err != nil {
		log.Error("failed to parse repository remote")
		return err
	}

	scm := detectSCM(rem.Host, ctx)
	if scm == git.Unrecognised {
		ctx.SCM = context.SCM{
			Provider: scm,
		}
		return nil
	}

	log.WithField("scm", scm).Info("scm provider identified")

	switch scm {
	case git.GitHub:
		ctx.SCM = github(rem)
	case git.GitLab:
		ctx.SCM = gitlab(rem)
	case git.CodeCommit:
		ctx.SCM = codecommit(rem)
	case git.Gitea:
		ctx.SCM = gitea(rem, ctx.Config.Gitea.URL)
	}

	return nil
}

func detectSCM(host string, ctx *context.Context) git.SCM {
	switch host {
	case "github.com":
		return git.GitHub
	case "gitlab.com":
		return git.GitLab
	}

	// Handle special case CodeCommit URLs
	if strings.HasPrefix(host, "git-codecommit") {
		return git.CodeCommit
	}

	if ctx.Config.GitHub != nil {
		if checkHost(host, ctx.Config.GitHub.URL) {
			return git.GitHub
		}
	} else if ctx.Config.GitLab != nil {
		if checkHost(host, ctx.Config.GitLab.URL) {
			return git.GitLab
		}
	} else if ctx.Config.Gitea != nil {
		if checkHost(host, ctx.Config.Gitea.URL) {
			return git.Gitea
		}
	}

	log.WithField("host", host).Warn("no recognised scm provider detected")
	return git.Unrecognised
}

func github(r git.Repository) context.SCM {
	url := fmt.Sprintf("https://%s/%s", r.Host, r.Path)

	return context.SCM{
		Provider:  git.GitHub,
		TagURL:    url + "/releases/tag/{{.Ref}}",
		CommitURL: url + "/commit/{{.Hash}}",
	}
}

func gitlab(r git.Repository) context.SCM {
	url := fmt.Sprintf("https://%s/%s", r.Host, r.Path)

	return context.SCM{
		Provider:  git.GitLab,
		TagURL:    url + "/-/tags/{{.Ref}}",
		CommitURL: url + "/-/commit/{{.Hash}}",
	}
}

func codecommit(r git.Repository) context.SCM {
	// CodeCommit URLs are a special case and require a region query parameter to be appended.
	// Extract the region from the clone URL
	t, _ := translate.RemoteHTTPS(r.Origin)

	// CodeCommit uses a different URL when browsing the repository
	browseURL := fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codecommit/repositories/%s", t.Region, r.Name)

	return context.SCM{
		Provider:  git.CodeCommit,
		TagURL:    fmt.Sprintf("%s/browse/refs/tags/{{.Ref}}?region=%s", browseURL, t.Region),
		CommitURL: fmt.Sprintf("%s/commit/{{.Hash}}?region=%s", browseURL, t.Region),
	}
}

func gitea(r git.Repository, u string) context.SCM {
	scheme := u[:strings.Index(u, ":")]
	url := fmt.Sprintf("%s://%s/%s", scheme, r.Host, r.Path)

	return context.SCM{
		Provider:  git.Gitea,
		TagURL:    url + "/releases/tag/{{.Ref}}",
		CommitURL: url + "/commit/{{.Hash}}",
	}
}

func checkHost(host string, in string) bool {
	u, err := url.Parse(in)
	if err != nil {
		log.WithField("url", in).Warn("could not parse provided URL")
		return false
	}

	return u.Host == host
}
