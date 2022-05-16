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

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelease(t *testing.T) {
	git.InitRepo(t)
	data := testFileWithConfig(t, "test.txt", ".uplift.yml")
	git.EmptyCommit(t, "feat: this is a release")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, tags[0].Ref, "0.1.0")

	// Ensure the tag is associated with the correct commit
	out, err := git.Clean(git.Run("tag", "-l", "0.1.0", `--format='%(subject)'`))
	require.NoError(t, err)
	assert.Contains(t, out, "ci(uplift):")

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.NotEqual(t, string(data), string(actual))
	assert.Contains(t, string(actual), "version: 0.1.0")

	assert.True(t, changelogExists(t))
}

func TestRelease_CheckFlag(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommits(t, "feat: this is a release", "Merge branch 'main' of https://github.com/test/repo")

	relCmd := newReleaseCmd(&globalOptions{}, os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--check"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)
}

func TestRelease_CheckFlagNoRelease(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "ci: not a release")

	relCmd := newReleaseCmd(&globalOptions{}, os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--check"})

	err := relCmd.Cmd.Execute()
	require.EqualError(t, err, "no release detected")
}

func TestRelease_PrereleaseFlag(t *testing.T) {
	git.InitRepo(t)
	testFileWithConfig(t, "test.txt", ".uplift.yml")
	git.EmptyCommit(t, "feat: this is a new feature")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--prerelease", "-beta.1+12345"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, "0.1.0-beta.1+12345", tags[0].Ref)

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Contains(t, string(actual), "version: 0.1.0-beta.1+12345")
}

func TestRelease_SkipChangelog(t *testing.T) {
	taggedRepo(t)

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--skip-changelog"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.False(t, changelogExists(t))
}

func TestRelease_SkipBumps(t *testing.T) {
	tagRepoWith(t, []string{"1.0.0"})
	testFileWithConfig(t, "test.txt", ".uplift.yml")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--skip-bumps"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.NotContains(t, string(actual), "version: 1.0.0")
}

func TestRelease_Hooks(t *testing.T) {
	git.InitRepo(t)
	configWithHooks(t)
	git.EmptyCommit(t, "feat: this is a new feature")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	require.Equal(t, 2, numHooksExecuted(t))
	assert.FileExists(t, BeforeFile)
	assert.FileExists(t, AfterFile)
}
