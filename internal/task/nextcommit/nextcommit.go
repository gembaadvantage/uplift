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

package nextcommit

import "context"

// Task ...
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "next-commit"
}

// Run ...
func (t Task) Run(ctx *context.Context) error {
	return nil
}

/*
func (b Bumper) buildCommit(ver string, commit git.CommitDetails) git.CommitDetails {
	c := git.CommitDetails{
		Author:  commit.Author,
		Email:   commit.Email,
		Message: fmt.Sprintf("ci(bump): bumped version to %s", ver),
	}

	if b.config.CommitAuthor.Name != "" {
		c.Author = b.config.CommitAuthor.Name
	}

	if b.config.CommitAuthor.Email != "" {
		c.Email = b.config.CommitAuthor.Email
	}

	if b.config.CommitMessage != "" {
		c.Message = b.config.CommitMessage
	}

	b.logger.Info("Any commits will use:\n%s", c)
	return c
}
*/
