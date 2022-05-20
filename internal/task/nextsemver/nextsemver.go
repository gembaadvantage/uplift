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

package nextsemver

import (
	semv "github.com/Masterminds/semver"
	"github.com/apex/log"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
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
	tag := git.LatestTag()
	if tag.Ref == "" {
		log.Debug("repository not tagged with version")
	}

	cl, err := git.Log(tag.Ref)
	if err != nil {
		return err
	}

	// Identify any commit that will trigger the largest semantic version bump
	inc := semver.ParseLog(cl)
	if inc == semver.NoIncrement {
		ctx.NoVersionChanged = true

		log.Warn("no commits trigger a change in semantic version")
		return nil
	}
	log.WithField("increment", string(inc)).Info("largest increment detected from commits")

	if tag.Ref == "" {
		// TODO: support a --no-prefix flag, that will strip the default 'v' prefix (only on first version)
		tag.Ref = "0.0.0"
	}

	pver, _ := semv.NewVersion(tag.Ref)

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
		}).Debug("appending prerelease version")
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

	log.WithField("version", ctx.NextVersion.Raw).Info("identified next semantic version")
	return nil
}
