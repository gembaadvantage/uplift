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
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// InitRepo creates an empty git repository within a temporary directory. Once created
// the current testing context will operate from within that directory until the calling
// test has completed
func InitRepo(t *testing.T) {
	MkTmpDir(t)

	// Initialise the git repo
	_, err := Run("init")
	require.NoError(t, err)

	EmptyCommit(t, "initialise repo")
}

// MkTmpDir creates an empty directory that is not a git repository. Once created the
// current testing context will operate from within that directory until the calling
// test has completed
func MkTmpDir(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	current, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(current))
	})
}

// EmptyCommit will create an empty commit without the need for modifying any existing files
// within the repository
func EmptyCommit(t *testing.T, commit string) {
	t.Helper()

	args := []string{
		"-c",
		"user.name='uplift'",
		"-c",
		"user.email='uplift@test.com'",
		"commit",
		"--allow-empty",
		"-m",
		commit,
	}

	_, err := Run(args...)
	require.NoError(t, err)
}

// EmptyCommits will create any number of empty commits without the need for modifying any
// existing files within the repository
func EmptyCommits(t *testing.T, commits ...string) {
	t.Helper()

	for _, msg := range commits {
		EmptyCommit(t, msg)
	}
}

// EmptyCommitAndTag will create an empty commit with an associated tag. No existing files
// will be modified within the repository
func EmptyCommitAndTag(t *testing.T, tag, msg string) {
	t.Helper()

	EmptyCommit(t, msg)
	err := Tag(tag)
	require.NoError(t, err)
}

// EmptyCommitsAndTag will create any number of empty commits and associate them with a tag.
// No existing files will be modified within the repository
func EmptyCommitsAndTag(t *testing.T, tag string, msgs ...string) {
	t.Helper()

	EmptyCommits(t, msgs...)
	err := Tag(tag)
	require.NoError(t, err)
}
