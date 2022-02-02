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
	"regexp"

	"github.com/apex/log"
	"github.com/gembaadvantage/codecommit-sign/pkg/translate"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
)

var (
	codecommitRgx = regexp.MustCompile(`^https://(.+@)?git-codecommit\.(.+)\.amazonaws.com/v1/repos/(.+)$`)
)

// Task that determines the SCM provider of a repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "detect scm"
}

// Skip is disabled for this task
func (t Task) Skip(ctx *context.Context) bool {
	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	rem, err := git.Remote()
	if err != nil {
		log.Info("failed to identify scm provider of repository")
		return err
	}

	ctx.SCM = &context.SCM{
		Provider: rem.Provider,
		URL:      rem.BrowseURL,
	}

	if rem.Provider == git.Unrecognised {
		log.Info("no recognised scm provider detected")
		return nil
	}

	log.WithField("scm", rem.Provider).Info("scm provider identified")

	// TODO: rename to template?

	switch rem.Provider {
	case git.GitHub:
		ctx.SCM.TagURL = rem.BrowseURL + "/releases/tag/{{.Ref}}"
		ctx.SCM.CommitURL = rem.BrowseURL + "/commit/{{.Hash}}"
	case git.GitLab:
		ctx.SCM.TagURL = rem.BrowseURL + "/-/tags/{{.Ref}}"
		ctx.SCM.CommitURL = rem.BrowseURL + "/-/commit/{{.Has}}"
	case git.CodeCommit:
		// CodeCommit URLs are a special case and require a region query parameter to be appended.
		// Extract the region from the clone URL
		t, _ := translate.RemoteHTTPS(rem.CloneURL)

		ctx.SCM.TagURL = fmt.Sprintf("%s/browse/refs/tags/{{.Ref}}?region=%s", rem.BrowseURL, t.Region)
		ctx.SCM.CommitURL = fmt.Sprintf("%s/commit/{{.Hash}}?region=%s", rem.BrowseURL, t.Region)
	}
	return nil
}
