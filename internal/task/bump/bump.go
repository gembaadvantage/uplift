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

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// FileBump defines how a version within a file will be matched through a regex
// and bumped using the provided version
type FileBump struct {
	Regex   string
	Version string
	Count   int
	SemVer  bool
}

// Task ...
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "bump"
}

// Run ...
func (t Task) Run(ctx *context.Context) error {
	if len(ctx.Config.Bumps) == 0 {
		return nil
	}

	n := 0
	for _, bump := range ctx.Config.Bumps {
		fb := FileBump{
			Regex:   bump.Regex,
			Version: ctx.NextVersion.Raw,
			Count:   bump.Count,
			SemVer:  bump.SemVer,
		}

		ok, err := bumpFile(ctx, bump.File, fb)
		if err != nil {
			return err
		}

		if ok {
			// Attempt to stage the changed file, if this isn't a dry run
			if !ctx.DryRun {
				if err := git.Stage(bump.File); err != nil {
					return err
				}
			}
			n++
		}
	}

	if n > 0 {
		return git.Commit(ctx.CommitDetails)
	}

	return nil
}

func bumpFile(ctx *context.Context, path string, bump FileBump) (bool, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		//b.logger.Warn("failed to open %s", path)
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
		//b.logger.Warn("version regex hasn't matched")
		return false, errors.New("no version matched in file")
	}
	mstr := string(m)

	if strings.Contains(mstr, bump.Version) {
		//b.logger.Info("skipped bumping %s as version already at %s", path, bump.Version)
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

	//b.logger.Success("bumped %s to version %s", path, bump.Version)

	// Don't make any file changes if part of a dry-run
	if ctx.DryRun {
		return false, nil
	}

	return true, ioutil.WriteFile(path, []byte(str), 0644)
}
