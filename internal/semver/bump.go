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
	"io/ioutil"
	"regexp"
	"strings"

	semv "github.com/Masterminds/semver"
	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/log"
)

const (
	firstVersion = "0.1.0"
)

// BumpOptions configures the behaviour when bumping a semantic version
type BumpOptions struct {
	Config  config.Uplift
	DryRun  bool
	Verbose bool
}

// Bumper is capable of bumping a semantic version associated with a git
// repository based on the conventional commits standard:
// @see https://www.conventionalcommits.org/en/v1.0.0/
type Bumper struct {
	logger log.ConsoleLogger
	config config.Uplift
	dryRun bool
}

// FileBump defines how a version within a file will be matched through a regex
// and bumped using the provided version
type FileBump struct {
	Regex   string
	Version string
	Count   int
}

var (
	version    = `v?\d+\.\d+\.\d+`
	versionRgx = regexp.MustCompile(version)
)

// NewBumper initialises a new semantic version bumper
func NewBumper(out io.Writer, opts BumpOptions) Bumper {
	l := log.NewSimpleLogger(out)
	if opts.Verbose {
		l = log.NewVerboseLogger(out)
	}

	// Override the first version if one hasn't been provided
	if opts.Config.FirstVersion == "" {
		opts.Config.FirstVersion = firstVersion
	}

	return Bumper{
		logger: l,
		config: opts.Config,
		dryRun: opts.DryRun,
	}
}

// Bump a semantic version based on the latest git log message within the associated
// git repository. Versions are incremented using the conventional commits standard.
// Once a version has been bumped, it will be tagged against the latest commit
func (b Bumper) Bump() error {
	if !git.IsRepo() {
		b.logger.Warn("no git repo found")
		return errors.New("current directory must be a git repo")
	}

	b.logger.Success("git repo found")

	meta, err := git.LatestCommit()
	if err != nil {
		b.logger.Warn("no commits found in repository")
		return err
	}
	b.logger.Success("retrieved latest commit:\n'%s'", meta.Message)

	inc := ParseCommit(meta.Message)
	if inc == NoIncrement {
		b.logger.Warn("commit doesn't contain a bump prefix, skipping!")
		return nil
	}
	b.logger.Success("commit contains a bump prefix, increment identified as '%s'", inc)

	ver := git.LatestTag()
	if ver == "" {
		ver = firstVersion
		b.logger.Success("no previous tags exist, using first version: %s\n", ver)
	} else {
		if ver, err = b.bumpVersion(ver, inc); err != nil {
			return err
		}
	}

	if err := b.bumpFiles(ver, meta); err != nil {
		return err
	}

	if b.dryRun {
		// Commit nothing on a dry run
		b.logger.Out(ver)
		return nil
	}

	if err := git.Tag(ver); err != nil {
		return err
	}

	b.logger.Out(ver)
	return nil
}

func (b Bumper) bumpVersion(v string, inc Increment) (string, error) {
	if inc == NoIncrement {
		return v, nil
	}

	b.logger.Info("existing version found: %s", v)

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
	case MajorIncrement:
		newVer = ver.IncMajor()
	case MinorIncrement:
		newVer = ver.IncMinor()
	case PatchIncrement:
		newVer = ver.IncPatch()
	}

	bv := fmt.Sprintf("%s%s", vp, newVer.String())
	b.logger.Success("bumped version to: %s", bv)

	return bv, nil
}

func (b Bumper) bumpFiles(v string, meta git.CommitMetadata) error {
	if len(b.config.Bumps) == 0 {
		b.logger.Info("no files to bump, skipping!")
		return nil
	}

	b.logger.Info("bumping files...")

	for _, bump := range b.config.Bumps {
		fb := FileBump{
			Regex:   bump.Regex,
			Version: v,
			Count:   bump.Count,
		}

		if err := b.bumpFile(bump.File, fb); err != nil {
			return err
		}

		if err := git.Stage(bump.File); err != nil {
			return err
		}

		b.logger.Success("bumped %s to version %s", bump.File, v)
	}

	// Don't commit anything
	if b.dryRun {
		return nil
	}

	// Commit all staged bumped files
	return git.Commit(meta.Author, meta.Email, "chore(release): release managed by uplift")
}

func (b Bumper) bumpFile(path string, bump FileBump) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		b.logger.Warn("failed to open %s", path)
		return err
	}

	// Ensure the supplied regex is valid, replacing the $VERSION token
	verRgx := strings.Replace(bump.Regex, "$VERSION", version, 1)

	rgx, err := regexp.Compile(verRgx)
	if err != nil {
		return err
	}

	m := rgx.Find(data)
	if m == nil {
		b.logger.Warn("version regex hasn't matched")
		return errors.New("no version matched in file")
	}

	// Use strings replace to ensure the replacement count is honoured
	n := -1
	if bump.Count > 0 {
		n = bump.Count
	}

	mstr := string(m)
	verRpl := versionRgx.ReplaceAllString(mstr, bump.Version)
	str := strings.Replace(string(data), mstr, verRpl, n)

	// Don't make any file changes if part of a dry-run
	if b.dryRun {
		return nil
	}

	return ioutil.WriteFile(path, []byte(str), 0644)
}
