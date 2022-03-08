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

func regexBump(ctx *context.Context, path string, bumps []config.RegexBump) (bool, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	str := string(data)

	for _, bump := range bumps {
		log.WithFields(log.Fields{
			"file":   path,
			"regex":  bump.Pattern,
			"count":  bump.Count,
			"semver": bump.SemVer,
		}).Debug("attempting file bump")

		mstr, err := match(bump.Pattern, str)
		if err != nil {
			return false, err
		}

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
		if bump.Count > 0 {
			n = bump.Count
		}

		// Strip any 'v' prefix if this must be a semantic version
		v := ctx.NextVersion.Raw
		if bump.SemVer {
			v = strictSemVer(v)
		}

		verRpl := semver.Regex.ReplaceAllString(mstr, v)
		str = strings.Replace(str, mstr, verRpl, n)
	}

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

func match(pattern string, data string) (string, error) {
	verRgx := strings.Replace(pattern, semver.Token, semver.Pattern, 1)

	rgx, err := regexp.Compile(verRgx)
	if err != nil {
		return "", err
	}

	m := rgx.FindString(data)
	if m == "" {
		return "", errors.New("no version matched in file")
	}

	return string(m), nil
}

func strictSemVer(v string) string {
	if v[0] == 'v' {
		log.Debug("trimming 'v' prefix from version")
		v = v[1:]
	}
	return v
}
