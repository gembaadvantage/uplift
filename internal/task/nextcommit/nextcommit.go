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

package nextcommit

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// Task for generating the next commit message
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "building next commit"
}

// Skip is disabled for this task
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.NoVersionChanged
}

// Run the task and generate the next commit by either impersonating the author
// from the last commit or by generating a user defined commit
func (t Task) Run(ctx *context.Context) error {
	c := git.CommitDetails{
		Author:  "uplift-bot",
		Email:   "uplift@gembaadvantage.com",
		Message: fmt.Sprintf("ci(uplift): uplifted for version %s", ctx.NextVersion.Raw),
	}

	if ctx.Config.CommitAuthor.Name != "" {
		log.Debug("overwriting commit author name from uplift config")
		c.Author = ctx.Config.CommitAuthor.Name
	}

	if ctx.Config.CommitAuthor.Email != "" {
		log.Debug("overwriting commit author email from uplift config")
		c.Email = ctx.Config.CommitAuthor.Email
	}

	if ctx.Config.CommitMessage != "" {
		log.Debug("overwriting commit message from uplift config")
		c.Message = strings.ReplaceAll(ctx.Config.CommitMessage, semver.Token, ctx.NextVersion.Raw)
	}

	author := git.Author()
	if author.Email != "" && author.Name != "" {
		log.Debug("overwriting commit author from git config")
		c.Author = author.Name
		c.Email = author.Email
	}

	ctx.CommitDetails = c
	log.WithFields(log.Fields{
		"name":    ctx.CommitDetails.Author,
		"email":   ctx.CommitDetails.Email,
		"message": ctx.CommitDetails.Message,
	}).Info("changes will be committed with")
	return nil
}
