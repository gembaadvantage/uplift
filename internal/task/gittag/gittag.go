/*
Copyright (c) 2021 Gemba Advantage

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

package gittag

import (
	"fmt"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
)

// Task for tagging a repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "git tag"
}

// Skip running the task if no version has changed
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.NoVersionChanged
}

// Run the task tagging a repository with the next semantic version. Supports both
// standard and annotated git tags
func (t Task) Run(ctx *context.Context) error {
	if ctx.CurrentVersion.Raw == ctx.NextVersion.Raw {
		log.WithFields(log.Fields{
			"current": ctx.CurrentVersion.Raw,
			"next":    ctx.NextVersion.Raw,
		}).Info("no version change detected")
		return nil
	}

	log.WithField("tag", ctx.NextVersion.Raw).Info("identified next tag")
	if ctx.DryRun {
		log.Info("skipping tag in dry run mode")
		return nil
	}

	if ctx.NextTagOnly {
		fmt.Fprint(ctx.Out, ctx.NextVersion.Raw)
		return nil
	}

	if ctx.Config.AnnotatedTags {
		log.Info("tagged repository with annotated tag")
		return git.AnnotatedTag(ctx.NextVersion.Raw, ctx.CommitDetails)
	}

	log.Info("tagged repository with standard tag")
	return git.Tag(ctx.NextVersion.Raw)
}
