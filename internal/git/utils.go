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
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Commit contains metadata about a specific git commit
type Commit struct {
	Message string `json:"message"`
	Author  string `json:"author"`
	Email   string `json:"email"`
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
	out, err := Clean(Run("log", "-1", `--pretty=format:'{"author": "%an", "email": "%ae", "message": "%B"}'`))
	if err != nil {
		return Commit{}, err
	}

	// Strip the final newline from the message and sanitise before unmarshalling
	out = strings.Replace(out, "\n\"}", "\"}", 1)
	out = strings.ReplaceAll(out, "\n", "\\n")

	var c Commit
	err = json.Unmarshal([]byte(out), &c)
	return c, err
}

// Tag the repository
func Tag(tag, author, email string) (string, error) {
	return Clean(Run(
		"-c",
		fmt.Sprintf("user.name='%s'", author),
		"-c",
		fmt.Sprintf("user.email='%s'", email),
		"tag",
		tag))
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
