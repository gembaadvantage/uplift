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
	"errors"
	"fmt"
	"io"

	semv "github.com/Masterminds/semver"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/log"
)

// BumpOptions ...
type BumpOptions struct {
	FirstVersion string
	DryRun       bool
	Verbose      bool
}

// Bumper ...
type Bumper struct {
	logger       log.ConsoleLogger
	firstVersion string
	dryRun       bool
}

// NewBumper ...
func NewBumper(out io.Writer, opts BumpOptions) Bumper {
	return Bumper{
		logger:       log.NewLogger(out, log.LoggerOptions{Debug: opts.Verbose}),
		firstVersion: opts.FirstVersion,
		dryRun:       opts.DryRun,
	}
}

// Bump ...
func (b Bumper) Bump() error {
	commit, err := git.LatestCommitMessage()
	if err != nil {
		return err
	}

	inc := checkIncrement(commit)
	if inc == noIncrement {
		return nil
	}

	ver, err := b.identifyVersion()
	if err != nil {
		return err
	}

	newVer, err := b.bumpVersion(ver, inc)
	if err != nil {
		return err
	}

	if b.dryRun {
		b.logger.Out("No changes will be committed to git repository")
		return nil
	}

	_, err = git.Tag(newVer)
	return err
}

func (b Bumper) identifyVersion() (string, error) {
	tag, err := git.LatestTag()
	if err != nil {
		if errors.Is(err, git.ErrNoTag) {
			return b.firstVersion, nil
		}
		return "", err
	}

	return tag, nil
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
