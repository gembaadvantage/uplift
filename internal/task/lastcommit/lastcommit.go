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
	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
)

// Task for reading the last commit message
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "inspect latest conventional commit"
}

// Skip is disabled for this task
func (t Task) Skip(ctx *context.Context) bool {
	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	// TODO: retrieve commit log between latest tag
	// TODO: scan messages trying to identify latest conventional commit
	// TODO: if none found, use latest commit

	cs, err := git.LatestCommits(ctx.CurrentVersion.Raw)
	if err != nil {
		log.Error("failed to retrieve latest commits")
		return err
	}

	commit, err := git.LatestCommit()
	if err != nil {
		log.Error("failed to retrieve latest commit")
		return err
	}
	log.WithFields(log.Fields{
		"author":  commit.Author,
		"email":   commit.Email,
		"message": commit.Message,
	}).Debug("retrieved latest commit")

	ctx.CommitDetails = commit
	return nil
}
