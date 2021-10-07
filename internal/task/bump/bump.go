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

package bump

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// FileBump defines how a version within a file will be matched through a regex
// and bumped using the provided version
type FileBump struct {
	Path    string
	Regex   string
	Version string
	Count   int
	SemVer  bool
}

// Task for bumping versions within files
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "bump"
}

// Skip running the task if no version has changed
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.NoVersionChanged
}

// Run the task bumping the semantic version of any file identified within
// the uplift configuration file
func (t Task) Run(ctx *context.Context) error {
	if len(ctx.Config.Bumps) == 0 {
		log.Info("no files to bump")
		return nil
	}

	n := 0
	for _, bump := range ctx.Config.Bumps {
		fb := FileBump{
			Path:    bump.File,
			Regex:   bump.Regex,
			Version: ctx.NextVersion.Raw,
			Count:   bump.Count,
			SemVer:  bump.SemVer,
		}

		ok, err := bumpFile(ctx, fb)
		if err != nil {
			return err
		}

		if ok {
			// Attempt to stage the changed file, if this isn't a dry run
			if !ctx.DryRun {
				if err := git.Stage(bump.File); err != nil {
					return err
				}
				log.WithField("file", bump.File).Info("successfully staged file")
			}
			n++
		}
	}

	if n > 0 {
		log.Info("attempting to commit file changes")
		return git.Commit(ctx.CommitDetails)
	}

	return nil
}

func bumpFile(ctx *context.Context, bump FileBump) (bool, error) {
	data, err := ioutil.ReadFile(bump.Path)
	if err != nil {
		return false, err
	}

	// Ensure the supplied regex is valid, replacing the $VERSION token
	verRgx := strings.Replace(bump.Regex, "$VERSION", semver.Pattern, 1)

	rgx, err := regexp.Compile(verRgx)
	if err != nil {
		return false, err
	}

	m := rgx.Find(data)
	if m == nil {
		return false, errors.New("no version matched in file")
	}
	mstr := string(m)

	if strings.Contains(mstr, bump.Version) {
		log.WithFields(log.Fields{
			"file":    bump.Path,
			"current": ctx.CurrentVersion.Raw,
			"next":    ctx.NextVersion.Raw,
		}).Info("skipping bump")
		return false, nil
	}

	// Use strings replace to ensure the replacement count is honoured
	n := -1
	if bump.Count > 0 {
		n = bump.Count
	}

	// Strip any 'v' prefix if this must be a semantic version
	v := bump.Version
	if bump.SemVer && v[0] == 'v' {
		v = v[1:]
	}

	verRpl := semver.Regex.ReplaceAllString(mstr, v)
	str := strings.Replace(string(data), mstr, verRpl, n)

	log.WithFields(log.Fields{
		"file":    bump.Path,
		"current": ctx.CurrentVersion.Raw,
		"next":    ctx.NextVersion.Raw,
	}).Info("file bumped")

	// Don't make any file changes if part of a dry-run
	if ctx.DryRun {
		log.Info("file not modified in dry run mode")
		return false, nil
	}

	return true, ioutil.WriteFile(bump.Path, []byte(str), 0644)
}
