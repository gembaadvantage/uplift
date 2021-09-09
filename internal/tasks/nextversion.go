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

package tasks

import (
	semv "github.com/Masterminds/semver"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// NextVersion ...
type NextVersion struct{}

// String ...
func (v NextVersion) String() string {
	return ""
}

// Run ...
func (v NextVersion) Run(ctx *context.Context) error {
	commit, err := git.LatestCommit()
	if err != nil {
		return err
	}

	inc := semver.ParseCommit(commit.Message)
	if inc == semver.NoIncrement {
		ctx.NextVersion = ctx.CurrentVersion
		return nil
	}

	pv, err := semv.NewVersion(ctx.CurrentVersion.Raw)
	if err != nil {
		return err
	}

	// Bump the semantic version based on the increment
	var nxt semv.Version
	switch inc {
	case semver.MajorIncrement:
		nxt = pv.IncMajor()
	case semver.MinorIncrement:
		nxt = pv.IncMinor()
	case semver.PatchIncrement:
		nxt = pv.IncPatch()
	}

	ctx.NextVersion = semver.Version{
		Prefix: ctx.CurrentVersion.Prefix,
		Patch:  uint64(nxt.Patch()),
		Minor:  uint64(nxt.Minor()),
		Major:  uint64(nxt.Major()),
		Raw:    nxt.Original(),
	}
	return nil
}
