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

package git

import (
	"errors"
	"os/exec"
	"strings"
)

// Commit contains metadata about a specific git commit
type Commit struct {
	Message string
	Author  string
	Email   string
}

// Run executes a git command and returns its output or errors
func Run(args ...string) (string, error) {
	var cmd = exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

// IsRepo identifies whether the current working directory is a recognised
// git repository
func IsRepo() bool {
	out, err := Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// LatestTag retrieves the latest tag within the repository
func LatestTag() string {
	tag, err := Clean(Run("describe", "--tags", "--abbrev=0"))
	if err != nil {
		return ""
	}

	return tag
}

// LatestCommit retrieves the latest commit within the repository
func LatestCommit() (Commit, error) {
	out, err := Clean(Run("log", "-1", `--pretty=format:'"%an","%ae","%B"'`))
	if err != nil {
		return Commit{}, err
	}

	// Split the formatted string into its component parts
	p := strings.Split(out, ",")

	// Strip quotes from around each part
	author := p[0][1 : len(p[0])-1]
	email := p[1][1 : len(p[1])-1]
	msg := p[2][1 : len(p[2])-1]

	return Commit{
		Author: author,
		Email:  email,
		// Strip trailing newline
		Message: strings.TrimRight(msg, "\n"),
	}, nil
}

// Tag the repository and push it to the origin
func Tag(tag string) error {
	if _, err := Clean(Run("tag", tag)); err != nil {
		return err
	}

	// Inspect the repo for an origin. If no origin exists, then skip the push
	if _, err := Clean(Run("remote", "show", "origin")); err != nil {
		return nil
	}

	if _, err := Clean(Run("push", "origin", tag)); err != nil {
		return err
	}

	return nil
}

// Clean the output
func Clean(output string, err error) (string, error) {
	// Preserve multi-line output, but trim the trailing newline
	output = strings.TrimSuffix(strings.Replace(output, "'", "", -1), "\n")
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}
