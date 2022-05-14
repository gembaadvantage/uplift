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

package lastcommit

import (
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// Task for reading the last commit message
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "scanning for conventional commit"
}

// Skip is disabled for this task
func (t Task) Skip(ctx *context.Context) bool {
	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	lc, err := git.LatestCommits(ctx.CurrentVersion.Raw)
	if err != nil {
		log.Error("failed to retrieve latest commits")
		return err
	}

	if len(lc) == 0 {
		log.Warn("no commits to scan, skipping...")
		return nil
	}

	// Default to the latest commit. If no other conventional commit
	// is found, this commit will be used by the remaining workflow
	ctx.CommitDetails = lc[0]

	for _, c := range lc {
		// Break as soon as a conventional commit is detected
		if semver.IsConventionalCommit(c.Message) {
			log.WithFields(log.Fields{
				"author":  c.Author,
				"email":   c.Email,
				"message": strings.TrimPrefix(c.Message, "\n"),
			}).Info("found commit")

			ctx.CommitDetails = c
			break
		}

		log.WithField("message", strings.TrimPrefix(c.Message, "\n")).Debug("skipping commit")
	}

	return nil
}
