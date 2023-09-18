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

package gitcheck

import (
	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
)

// Task for detecting if uplift is being run within a recognised
// git repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "checking git"
}

// Skip running the task
func (t Task) Skip(_ *context.Context) bool {
	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	repo, err := ctx.GitClient.Repository()
	if err != nil {
		return ErrNoRepository
	}

	log.Debug("checking if repository is dirty")
	status, err := ctx.GitClient.PorcelainStatus()
	if err != nil {
		return err
	}

	if len(status) > 0 {

		if len(ctx.IncludeArtifacts) == 0 {
			return ErrDirty{status}
		}

		for _, sts := range status {
			if !stringInSlice(sts.Path, ctx.IncludeArtifacts) {
				return ErrDirty{status}
			}
		}
	}

	log.Debug("checking for detached head")
	if repo.DetachedHead {
		if ctx.IgnoreDetached {
			log.Warn("detached HEAD detected. This may impact certain operations within uplift")
		} else {
			return ErrDetachedHead{}
		}
	}

	log.Debug("checking if shallow clone used")
	if repo.ShallowClone {
		if ctx.IgnoreShallow {
			log.Warn("shallow clone detected. This may impact certain operations within uplift")
		} else {
			return ErrShallowClone{}
		}
	}

	return nil
}

func stringInSlice(stringToFind string, listOfStrings []string) bool {
	for _, value := range listOfStrings {
		if value == stringToFind {
			return true
		}
	}
	return false
}
