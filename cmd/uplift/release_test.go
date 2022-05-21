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
	untaggedRepo(t, "ci: update pipeline", "docs: update docs", "fix: bug fix", "feat: new feature")
	data := testFileWithConfig(t, "test.txt", ".uplift.yml")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, tags[0].Ref, "v0.1.0")

	// Ensure the tag is associated with the correct commit
	out, err := git.Clean(git.Run("tag", "-l", "v0.1.0", `--format='%(subject)'`))
	require.NoError(t, err)
	assert.Equal(t, out, "ci(uplift): uplifted for version v0.1.0")

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.NotEqual(t, string(data), string(actual))
	assert.Contains(t, string(actual), "version: v0.1.0")

	assert.True(t, changelogExists(t))
	cl := readChangelog(t)
	assert.Contains(t, cl, "## v0.1.0")
}

func TestRelease_NoPrefix(t *testing.T) {
	untaggedRepo(t,
		"ci: update pipeline",
		"docs: update docs",
		`refactor: a big change
a description about the work involved

BREAKING CHANGE: the existing cli is no longer backward compatible`,
		"fix: bug fix",
		"feat: new feature")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--no-prefix"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, tags[0].Ref, "1.0.0")
}

func TestRelease_CheckFlag(t *testing.T) {
	untaggedRepo(t, "Merge branch 'main' of https://github.com/test/repo", "feat: new feature", "docs: update docs", "ci: workflow")

	relCmd := newReleaseCmd(&globalOptions{}, os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--check"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)
}

func TestRelease_CheckFlagNoRelease(t *testing.T) {
	untaggedRepo(t, "ci: not a release", "docs: update docs", "refactor: change everything")

	relCmd := newReleaseCmd(&globalOptions{}, os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--check"})

	err := relCmd.Cmd.Execute()
	require.EqualError(t, err, "no release detected")
}

func TestRelease_PrereleaseFlag(t *testing.T) {
	untaggedRepo(t, "docs: update docs", "feat: new feature", "refactor: make changes")
	testFileWithConfig(t, "test.txt", ".uplift.yml")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--prerelease", "-beta.1+12345"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0-beta.1+12345", tags[0].Ref)

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Contains(t, string(actual), "version: v0.1.0-beta.1+12345")
}

func TestRelease_SkipChangelog(t *testing.T) {
	taggedRepo(t, "1.0.0", "feat: first feature")

	// Ensure another release would be triggered
	git.EmptyCommits(t, "ci: updated workflow", "fix: bug fix", "docs: updated docs")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--skip-changelog"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tag := git.LatestTag()
	assert.Equal(t, "1.0.1", tag.Ref)

	assert.False(t, changelogExists(t))
}

func TestRelease_SkipBumps(t *testing.T) {
	taggedRepo(t, "1.0.0", "feat: first feature")
	testFileWithConfig(t, "test.txt", ".uplift.yml")

	// Ensure another release would be triggered
	git.EmptyCommits(t, "ci: updated workflow", "fix: bug fix", "docs: updated docs")

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--skip-bumps"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tag := git.LatestTag()
	assert.Equal(t, "1.0.1", tag.Ref)

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.NotContains(t, string(actual), "version: 1.0.0")
}

func TestRelease_Hooks(t *testing.T) {
	untaggedRepo(t, "docs: updated docs", "feat: new feature", "ci: update workflow")
	configWithHooks(t)

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	require.Equal(t, 8, numHooksExecuted(t))
	assert.FileExists(t, BeforeFile)
	assert.FileExists(t, BeforeBumpFile)
	assert.FileExists(t, AfterBumpFile)
	assert.FileExists(t, BeforeChangelogFile)
	assert.FileExists(t, AfterChangelogFile)
	assert.FileExists(t, BeforeTagFile)
	assert.FileExists(t, AfterTagFile)
	assert.FileExists(t, AfterFile)
}
