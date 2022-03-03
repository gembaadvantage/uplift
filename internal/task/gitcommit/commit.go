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

package gitcommit

import (
	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
)

// Task for committing all staged changes and optionally pushing
// them to a git remote
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "git commit"
}

// Skip running the task
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.DryRun || ctx.NoVersionChanged
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	stg, err := git.Staged()
	if err != nil {
		return err
	}

	if len(stg) == 0 {
		log.Info("no files staged skipping push")
		return nil
	}

	log.Info("commit outstanding changes")
	if err := git.Commit(ctx.CommitDetails); err != nil {
		return err
	}

	if ctx.NoPush {
		log.Info("skipping push of commit to remote")
		return nil
	}

	log.Info("pushing commit to remote")
	return git.Push()
}
