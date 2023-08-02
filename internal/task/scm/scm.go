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
)

type remote struct {
	Origin string
	Owner  string
	Name   string
	Host   string
	Path   string
}

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
	repo, err := ctx.GitClient.Repository()
	if err != nil {
		log.Error("failed to parse repository remote")
		return err
	}

	rem, err := parseRemote(repo.Origin)
	if err != nil {
		return err
	}

	scm := detectSCM(rem.Host, ctx)
	if scm == context.Unrecognised {
		ctx.SCM = context.SCM{
			Provider: scm,
		}
		return nil
	}

	log.WithField("scm", scm).Info("scm provider identified")

	switch scm {
	case context.GitHub:
		ctx.SCM = github(rem)
	case context.GitLab:
		ctx.SCM = gitlab(rem)
	case context.CodeCommit:
		ctx.SCM = codecommit(rem)
	case context.Gitea:
		ctx.SCM = gitea(rem, ctx.Config.Gitea.URL)
	}

	return nil
}

func parseRemote(remURL string) (remote, error) {
	origin := remURL

	// Strip off any trailing .git suffix
	rem := strings.TrimSuffix(remURL, ".git")

	// Special use case for CodeCommit as that prefixes SSH URLs with ssh://
	rem = strings.TrimPrefix(rem, "ssh://")

	// Detect and translate a CodeCommit GRC URL into its HTTPS counterpart
	if strings.HasPrefix(rem, "codecommit:") {
		// Translate a codecommit GRC URL into its HTTPS counterpart
		var err error
		if rem, err = translate.FromGRC(rem); err != nil {
			return remote{}, err
		}

		origin = rem
	}

	if strings.HasPrefix(rem, "git@") {
		// Sanitise any SSH based URL to ensure it is parseable
		rem = strings.TrimPrefix(rem, "git@")
		rem = strings.Replace(rem, ":", "/", 1)
	}

	u, err := url.Parse(rem)
	if err != nil {
		return remote{}, err
	}

	// Split into parts
	p := strings.Split(u.Path, "/")

	if len(p) < 3 {
		// This could be a custom Git server that doesn't follow the expected pattern.
		// Don't fail, but return the raw origin for custom parsing
		return remote{Origin: origin}, nil
	}

	// If the repository has a HTTP(S) origin, the host will have been correctly identified
	host := u.Host
	if host == "" {
		// For other schemes, assume the host is contained within the first part
		host = p[0]
	}
	owner := p[1]
	name := p[len(p)-1]
	path := strings.Join(p[1:], "/")

	if strings.Contains(host, "codecommit") {
		// No concept of an owner with CodeCommit repositories
		owner = ""
		path = name
	}

	return remote{
		Origin: origin,
		Owner:  owner,
		Name:   name,
		Path:   path,
		Host:   host,
	}, nil
}

func detectSCM(host string, ctx *context.Context) context.SCMProvider {
	switch host {
	case "github.com":
		return context.GitHub
	case "gitlab.com":
		return context.GitLab
	}

	// Handle special case CodeCommit URLs
	if strings.HasPrefix(host, "git-codecommit") {
		return context.CodeCommit
	}

	if ctx.Config.GitHub != nil {
		if checkHost(host, ctx.Config.GitHub.URL) {
			return context.GitHub
		}
	} else if ctx.Config.GitLab != nil {
		if checkHost(host, ctx.Config.GitLab.URL) {
			return context.GitLab
		}
	} else if ctx.Config.Gitea != nil {
		if checkHost(host, ctx.Config.Gitea.URL) {
			return context.Gitea
		}
	}

	log.WithField("host", host).Warn("no recognised scm provider detected")
	return context.Unrecognised
}

func github(rem remote) context.SCM {
	url := fmt.Sprintf("https://%s/%s", rem.Host, rem.Path)

	return context.SCM{
		Provider:  context.GitHub,
		TagURL:    url + "/releases/tag/{{.Ref}}",
		CommitURL: url + "/commit/{{.Hash}}",
	}
}

func gitlab(rem remote) context.SCM {
	url := fmt.Sprintf("https://%s/%s", rem.Host, rem.Path)

	return context.SCM{
		Provider:  context.GitLab,
		TagURL:    url + "/-/tags/{{.Ref}}",
		CommitURL: url + "/-/commit/{{.Hash}}",
	}
}

func codecommit(rem remote) context.SCM {
	// CodeCommit URLs are a special case and require a region query parameter to be appended.
	// Extract the region from the clone URL
	t, _ := translate.RemoteHTTPS(rem.Origin)

	// CodeCommit uses a different URL when browsing the repository
	browseURL := fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codecommit/repositories/%s", t.Region, rem.Name)

	return context.SCM{
		Provider:  context.CodeCommit,
		TagURL:    fmt.Sprintf("%s/browse/refs/tags/{{.Ref}}?region=%s", browseURL, t.Region),
		CommitURL: fmt.Sprintf("%s/commit/{{.Hash}}?region=%s", browseURL, t.Region),
	}
}

func gitea(rem remote, u string) context.SCM {
	scheme := u[:strings.Index(u, ":")]
	url := fmt.Sprintf("%s://%s/%s", scheme, rem.Host, rem.Path)

	return context.SCM{
		Provider:  context.Gitea,
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
