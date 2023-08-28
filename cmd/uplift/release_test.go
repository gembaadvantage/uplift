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
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelease(t *testing.T) {
	log := `feat: new feature
fix: bug fix
docs: update docs
ci: update pipeline`
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithCommittedFiles("test.txt", ".uplift.yml"),
		gittest.WithFileContent("test.txt", bumpFile, ".uplift.yml", bumpConfig))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	assert.Len(t, tags, 1)
	assert.Equal(t, tags[0], "v0.1.0")

	// Ensure the tag is associated with the correct commit
	out := gittest.MustExec(t, `git tag -l v0.1.0 --format='%(subject)'`)
	require.NoError(t, err)
	assert.Equal(t, out, "ci(uplift): uplifted for version v0.1.0")

	actual, err := os.ReadFile("test.txt")
	require.NoError(t, err)
	assert.NotEqual(t, string(bumpFile), string(actual))
	assert.Contains(t, string(actual), "version: v0.1.0")

	assert.True(t, changelogExists(t))
	cl := readChangelog(t)
	assert.Contains(t, cl, "## v0.1.0")
}

func TestRelease_NoPrefix(t *testing.T) {
	log := `> ci: update pipeline
> docs: update docs
> refactor: a big change
a description about the work involved
BREAKING CHANGE: the existing cli is no longer backward compatible
> fix: bug fix
> feat: new feature`
	gittest.InitRepository(t, gittest.WithLog(log))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--no-prefix"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.Len(t, tags, 1)
	assert.Equal(t, "1.0.0", tags[0])
}

func TestRelease_CheckFlag(t *testing.T) {
	log := `ci: workflow
docs: update docs
feat: new feature
Merge branch 'main' of https://github.com/test/repo`
	gittest.InitRepository(t, gittest.WithLog(log))

	relCmd := newReleaseCmd(&globalOptions{}, os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--check"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)
}

func TestRelease_CheckFlagNoRelease(t *testing.T) {
	log := `refactor: change everything
docs: update docs
ci: not a release`
	gittest.InitRepository(t, gittest.WithLog(log))

	relCmd := newReleaseCmd(&globalOptions{}, os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--check"})

	err := relCmd.Cmd.Execute()
	require.EqualError(t, err, "no release detected")
}

func TestRelease_PrereleaseFlag(t *testing.T) {
	log := `refactor: make changes
feat: new feature
docs: update docs`
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithCommittedFiles("test.txt", ".uplift.yml"),
		gittest.WithFileContent("test.txt", bumpFile, ".uplift.yml", bumpConfig))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--prerelease", "-beta.1+12345"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	require.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0-beta.1+12345", tags[0])

	actual, err := os.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Contains(t, string(actual), "version: v0.1.0-beta.1+12345")
}

func TestRelease_SkipChangelog(t *testing.T) {
	log := `docs: updated docs
fix: bug fix
ci: updated workflow
(tag: 1.0.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithLog(log))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--skip-changelog"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	assert.Contains(t, tags, "1.0.1")
	assert.False(t, changelogExists(t))
}

func TestRelease_SkipBumps(t *testing.T) {
	log := `docs: updated docs
fix: bug fix
ci: updated workflow
(tag: 1.0.0) feat: first feature`
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithCommittedFiles("test.txt", ".uplift.yml"),
		gittest.WithFileContent("test.txt", bumpFile, ".uplift.yml", bumpConfig))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--skip-bumps"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	assert.Contains(t, tags, "1.0.1")

	actual, err := os.ReadFile("test.txt")
	require.NoError(t, err)
	assert.NotContains(t, string(actual), "version: 1.0.1")
}

func TestRelease_Hooks(t *testing.T) {
	log := `ci: update workflow
feat: new feature
docs: updated docs`
	gittest.InitRepository(t, gittest.WithLog(log))
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

func TestRelease_ExcludesUpliftCommitByDefault(t *testing.T) {
	log := `fix: a bug fix
ci: tweak workflow`
	gittest.InitRepository(t, gittest.WithLog(log))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.NotContains(t, cl, "ci(uplift): uplifted version")
	assert.Contains(t, cl, "fix: a bug fix")
	assert.Contains(t, cl, "ci: tweak workflow")
}

func TestRelease_WithExclude(t *testing.T) {
	log := `docs: some new docs
ci: a ci task
fix: a new fix
feat: a new feat`
	gittest.InitRepository(t, gittest.WithLog(log))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--exclude", "^ci,^docs"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, "feat: a new feat")
	assert.Contains(t, cl, "fix: a new fix")
	assert.NotContains(t, cl, "ci: a ci task")
	assert.NotContains(t, cl, "docs: some new docs")
}

func TestRelease_WithInclude(t *testing.T) {
	log := `docs: some new docs
ci: a ci task
fix: a new fix
feat: a new feat`
	gittest.InitRepository(t, gittest.WithLog(log))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--include", "^feat"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, "feat: a new feat")
	assert.NotContains(t, cl, "fix: a new fix")
	assert.NotContains(t, cl, "ci: a ci task")
	assert.NotContains(t, cl, "docs: some new docs")
}

func TestRelease_WithMultiline(t *testing.T) {
	log := `> feat: this is a multiline commit
The entire contents of this commit should exist in the changelog.

Multiline formatting should be correct for rendering in markdown
> fix: this is a bug fix
> docs: update documentation
this now includes code examples`
	gittest.InitRepository(t, gittest.WithLog(log))

	relCmd := newReleaseCmd(noChangesPushed(), os.Stdout)
	relCmd.Cmd.SetArgs([]string{"--multiline"})

	err := relCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, `feat: this is a multiline commit
  The entire contents of this commit should exist in the changelog.

  Multiline formatting should be correct for rendering in markdown`)
	assert.Contains(t, cl, "fix: this is a bug fix")
	assert.Contains(t, cl, `docs: update documentation
  this now includes code examples`)
}
