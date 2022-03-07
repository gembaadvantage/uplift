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

package bump

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// TODO: support iterating over slice and bump using all regex patterns

func regexBump(ctx *context.Context, path string, bumps []config.RegexBump) (bool, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	// Ensure the supplied regex is valid, replacing the $VERSION token
	verRgx := strings.Replace(bumps[0].Pattern, semver.Token, semver.Pattern, 1)

	rgx, err := regexp.Compile(verRgx)
	if err != nil {
		return false, err
	}

	m := rgx.Find(data)
	if m == nil {
		return false, errors.New("no version matched in file")
	}
	mstr := string(m)

	if strings.Contains(mstr, ctx.NextVersion.Raw) {
		log.WithFields(log.Fields{
			"file":    path,
			"current": ctx.CurrentVersion.Raw,
			"next":    ctx.NextVersion.Raw,
		}).Info("skipping bump")
		return false, nil
	}

	// Use strings replace to ensure the replacement count is honoured
	n := -1
	if bumps[0].Count > 0 {
		n = bumps[0].Count
	}

	// Strip any 'v' prefix if this must be a semantic version
	v := ctx.NextVersion.Raw
	if bumps[0].SemVer && v[0] == 'v' {
		v = v[1:]
	}

	verRpl := semver.Regex.ReplaceAllString(mstr, v)
	str := strings.Replace(string(data), mstr, verRpl, n)

	log.WithFields(log.Fields{
		"file":    path,
		"current": ctx.CurrentVersion.Raw,
		"next":    ctx.NextVersion.Raw,
	}).Info("file bumped")

	// Don't make any file changes if part of a dry-run
	if ctx.DryRun {
		log.Info("file not modified in dry run mode")
		return false, nil
	}

	return true, ioutil.WriteFile(path, []byte(str), 0644)
}
