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

package context

import (
	ctx "context"
	"io"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// Context provides a way to share common state across tasks
type Context struct {
	ctx.Context
	Out              io.Writer
	Config           config.Uplift
	DryRun           bool
	Debug            bool
	CurrentVersion   semver.Version
	NextVersion      semver.Version
	NoVersionChanged bool
	CommitDetails    git.CommitDetails
	FetchTags        bool
	NextTagOnly      bool
	NoPush           bool
}

// New constructs a context that captures both runtime configuration and
// user defined runtime options
func New(cfg config.Uplift, out io.Writer) *Context {
	return &Context{
		Context: ctx.Background(),
		Config:  cfg,
		Out:     out,
	}
}
