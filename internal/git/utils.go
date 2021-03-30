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

var (
	// ErrNotTag is thrown if a repository does not contain any tags
	ErrNoTag = errors.New("no tag exists in repository")
)

// Run executes a git command and returns its output or errors
func Run(args ...string) (string, error) {
	var cmd = exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

// LatestTag retrieves the latest tag within the repository
func LatestTag() (string, error) {
	out, err := Clean(Run("describe", "--tags", "--abbrev=0"))
	if err != nil {
		if strings.Contains(err.Error(), "No names found, cannot describe anything") {
			return out, ErrNoTag
		}
	}

	return out, err
}

// LatestCommitMessage retrieves the latest commit message within the repository
func LatestCommitMessage() (string, error) {
	return Clean(Run("log", "-1", "--pretty=format:%B"))
}

// Tag the repository
func Tag(tag string) (string, error) {
	return Clean(Run("tag", tag))
}

// Clean the output
func Clean(output string, err error) (string, error) {
	output = strings.Replace(strings.Split(output, "\n")[0], "'", "", -1)
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}
