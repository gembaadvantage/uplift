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

package semver

import (
	"fmt"
	"io"

	semv "github.com/Masterminds/semver"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/log"
)

// BumpOptions configures the behaviour when bumping a semantic version
type BumpOptions struct {
	FirstVersion string
	DryRun       bool
	Verbose      bool
}

// Bumper is capable of bumping a semantic version associated with a git
// repository based on the conventional commits standard:
// @see https://www.conventionalcommits.org/en/v1.0.0/
type Bumper struct {
	logger       log.ConsoleLogger
	firstVersion string
	dryRun       bool
}

// NewBumper initialises a new semantic version bumper
func NewBumper(out io.Writer, opts BumpOptions) Bumper {
	return Bumper{
		logger:       log.NewLogger(out, log.LoggerOptions{Debug: opts.Verbose}),
		firstVersion: opts.FirstVersion,
		dryRun:       opts.DryRun,
	}
}

// Bump a semantic version based on the latest git log message within the associated
// git repository. Versions are incremented using the conventional commits standard.
// Once a version has been bumped, it will be tagged against the latest commit
func (b Bumper) Bump() error {
	if !git.IsRepo() {
		// TODO: log
		return nil
	}

	commit, err := git.LatestCommitMessage()
	if err != nil {
		return err
	}

	inc := ParseCommit(commit)
	if inc == noIncrement {
		// TODO: no version was bumped, skip
		return nil
	}

	ver := git.LatestTag()
	if ver == "" {
		ver = b.firstVersion
	} else {
		if ver, err = b.bumpVersion(ver, inc); err != nil {
			return err
		}
	}

	if b.dryRun {
		b.logger.Out("No changes will be committed to git repository")
		return nil
	}

	_, err = git.Tag(ver)
	return err
}

func (b Bumper) bumpVersion(v string, inc increment) (string, error) {
	if inc == noIncrement {
		return v, nil
	}

	ver, err := semv.NewVersion(v)
	if err != nil {
		return "", err
	}

	// If the provided version has a "v" prefix, ensure it is preserved in the new version
	vp := ""
	if v[0] == 'v' {
		vp = "v"
	}

	var newVer semv.Version

	switch inc {
	case majorIncrement:
		newVer = ver.IncMajor()
	case minorIncrement:
		newVer = ver.IncMinor()
	case patchIncrement:
		newVer = ver.IncPatch()
	}

	return fmt.Sprintf("%s%s", vp, newVer.String()), nil
}
