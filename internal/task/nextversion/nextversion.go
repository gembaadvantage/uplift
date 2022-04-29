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

package nextversion

import (
	semv "github.com/Masterminds/semver"
	"github.com/apex/log"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
)

const (
	defaultVersion = "0.1.0"
)

// Task that determines the next semantic version of a repository
// based on the conventional commit of the last commit message
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "next semantic version"
}

// Skip is disabled for this task
func (t Task) Skip(ctx *context.Context) bool {
	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	inc := semver.ParseCommit(ctx.CommitDetails.Message)
	if inc == semver.NoIncrement {
		ctx.NextVersion = ctx.CurrentVersion
		ctx.NoVersionChanged = true
		log.WithField("commit", ctx.CommitDetails.Message).Warn("commit doesn't trigger change in semantic version")
		return nil
	}

	log.WithField("increment", string(inc)).Debug("increment detected from commit")

	// If this is the first tag, use the required default
	if ctx.CurrentVersion.Raw == "" {
		ctx.NextVersion, _ = semver.Parse(firstVersion(ctx))
		log.WithField("version", ctx.NextVersion.Raw).Info("repository not tagged, using first version")
		return nil
	}

	pver, _ := semv.NewVersion(ctx.CurrentVersion.Raw)

	// Bump the semantic version based on the increment
	var nxt semv.Version
	switch inc {
	case semver.MajorIncrement:
		nxt = pver.IncMajor()
	case semver.MinorIncrement:
		nxt = pver.IncMinor()
	case semver.PatchIncrement:
		nxt = pver.IncPatch()
	}

	// Append any prerelease suffixes
	if ctx.Prerelease != "" {
		nxt, _ = nxt.SetPrerelease(ctx.Prerelease)
		nxt, _ = nxt.SetMetadata(ctx.Metadata)

		log.WithFields(log.Fields{
			"prerelease": ctx.Prerelease,
			"metadata":   ctx.Metadata,
		}).Info("appending prerelease version")
	}

	ctx.NextVersion = semver.Version{
		Prefix:     ctx.CurrentVersion.Prefix,
		Patch:      nxt.Patch(),
		Minor:      nxt.Minor(),
		Major:      nxt.Major(),
		Prerelease: nxt.Prerelease(),
		Metadata:   nxt.Metadata(),
		Raw:        nxt.Original(),
	}

	log.WithField("next", ctx.NextVersion.Raw).Info("identified next semantic version")
	return nil
}

func firstVersion(ctx *context.Context) string {
	fv := defaultVersion

	if ctx.Config.FirstVersion != "" {
		fv = ctx.Config.FirstVersion
		log.WithField("version", fv).Debug("setting first version based on config")
	}

	// Append any semantic prerelease suffix if needed
	v, _ := semv.NewVersion(fv)
	*v, _ = v.SetPrerelease(ctx.Prerelease)
	*v, _ = v.SetMetadata(ctx.Metadata)

	return v.Original()
}
