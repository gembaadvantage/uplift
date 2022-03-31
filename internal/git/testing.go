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

package git

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	InitCommit = "initialise repo"
)

// InitRepo creates an empty git repository within a temporary directory. Once created
// the current testing context will operate from within that directory until the calling
// test has completed
func InitRepo(t *testing.T) string {
	t.Helper()

	MkTmpDir(t)

	// Initialise the git repo
	_, err := Run("init")
	require.NoError(t, err)

	// Set a default origin that can be changed if needed
	RemoteOrigin(t, "http://example.com/project/repository")

	return EmptyCommit(t, InitCommit)
}

// RemoteOrigin sets the URL of the remote origin associated with the current git repository
func RemoteOrigin(t *testing.T, url string) {
	t.Helper()

	// If an origin already exists ensure it is removed first
	Run("remote", "remove", "origin")

	_, err := Run("remote", "add", "origin", url)
	require.NoError(t, err)

	_, err = Run("config", "branch.main.remote", "origin")
	require.NoError(t, err)

	_, err = Run("config", "branch.main.merge", "refs/heads/main")
	require.NoError(t, err)
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
func EmptyCommit(t *testing.T, commit string) string {
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

	// Grab the unabbreviated hash of the newly created commit
	out, err := Clean(Run("rev-parse", "HEAD"))
	require.NoError(t, err)

	return out
}

// EmptyCommits will create any number of empty commits without the need for modifying any
// existing files within the repository
func EmptyCommits(t *testing.T, commits ...string) []string {
	t.Helper()

	hs := make([]string, 0, len(commits))
	for _, msg := range commits {
		hs = append(hs, EmptyCommit(t, msg))
	}
	return hs
}

// EmptyCommitAndTag will create an empty commit with an associated tag. No existing files
// will be modified within the repository
func EmptyCommitAndTag(t *testing.T, tag, msg string) string {
	t.Helper()

	h := EmptyCommit(t, msg)
	err := Tag(tag)
	require.NoError(t, err)

	return h
}

// EmptyCommitsAndTag will create any number of empty commits and associate them with a tag.
// No existing files will be modified within the repository
func EmptyCommitsAndTag(t *testing.T, tag string, msgs ...string) []string {
	t.Helper()

	hs := EmptyCommits(t, msgs...)
	err := Tag(tag)
	require.NoError(t, err)

	return hs
}

// TimeBasedTagSeries is a specialised utility function for generating a series of tags
// that are spaced apart by a day. This is important for any tests that require
// time based filtering of tags. If all tags are created at the same time, filtering
// can produce inconsistent ordering. Commits are auto-generated with the following
// format:
//
// feat: <COMMIT_INDEX>
//
// e.g. []tags{"1.0.0", "2.0.0"} => feat: 1, feat: 2
func TimeBasedTagSeries(t *testing.T, tags []string) []string {
	// Ensure the GIT_COMMITTER_DATE is always reset
	defer func() {
		os.Unsetenv("GIT_COMMITTER_DATE")
	}()

	h := make([]string, 0, len(tags))

	// Calculate the max days in the past
	max := len(tags)

	now := time.Now().UTC()
	for i, c := 0, 1; i < len(tags); i, c = i+1, c+1 {
		dt := now.AddDate(0, 0, -(max - i)).Format(time.RFC3339)

		fmt.Println(dt)

		// Based on the git spec, this env var should be set when manipulating dates
		// of tags and commits
		os.Setenv("GIT_COMMITTER_DATE", dt)

		args := []string{
			"-c",
			"user.name='uplift'",
			"-c",
			"user.email='uplift@test.com'",
			"commit",
			"--allow-empty",
			"-m",
			fmt.Sprintf("feat: %d", c),
			"--date",
			dt,
		}

		_, err := Run(args...)
		require.NoError(t, err)

		// Grab the unabbreviated hash of the newly created commit
		out, err := Clean(Run("rev-parse", "HEAD"))
		require.NoError(t, err)

		h = append(h, out)

		// Ensure the tag is generated with the same date
		err = Tag(tags[i])
		require.NoError(t, err)
	}

	return h
}
