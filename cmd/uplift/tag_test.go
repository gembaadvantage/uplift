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
	"bytes"
	"os"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	untaggedRepo(t, "fix: bug fix", "docs: update docs", "ci: update pipeline", "feat: new feature")

	tagCmd := newTagCmd(noChangesPushed(), os.Stdout)

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0", tags[0].Ref)
}

func TestTag_CurrentFlag(t *testing.T) {
	taggedRepo(t, "v0.1.0", "docs: updated docs", "fix: bug fix")

	var buf bytes.Buffer
	tagCmd := newTagCmd(noChangesPushed(), &buf)
	tagCmd.Cmd.SetArgs([]string{"--current"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0", buf.String())
}

func TestTag_NextFlag(t *testing.T) {
	untaggedRepo(t, "docs: updated docs", "refactor!: breaking cli change", "ci: update pipeline", "fix: bug fix")

	var buf bytes.Buffer
	tagCmd := newTagCmd(noChangesPushed(), &buf)
	tagCmd.Cmd.SetArgs([]string{"--next"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 0)
	assert.Equal(t, "v1.0.0", buf.String())
}

func TestTag_CurrentAndNextFlag(t *testing.T) {
	taggedRepo(t, "v0.1.0", "docs: updated docs", "fix: bug fix")
	git.EmptyCommit(t, "fix: found another bug")

	var buf bytes.Buffer
	tagCmd := newTagCmd(noChangesPushed(), &buf)
	tagCmd.Cmd.SetArgs([]string{"--current", "--next"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0 v0.1.1", buf.String())
}

func TestTag_NoPrefix(t *testing.T) {
	untaggedRepo(t, "docs: update docs", "fix: bug fix")

	tagCmd := newTagCmd(noChangesPushed(), os.Stdout)
	tagCmd.Cmd.SetArgs([]string{"--no-prefix"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, "0.0.1", tags[0].Ref)
}

func TestTag_PrereleaseFlag(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: a new feature")

	tagCmd := newTagCmd(noChangesPushed(), os.Stdout)
	tagCmd.Cmd.SetArgs([]string{"--prerelease", "-beta.1+12345"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0-beta.1+12345", tags[0].Ref)
}

func TestTag_Hooks(t *testing.T) {
	git.InitRepo(t)
	configWithHooks(t)
	git.EmptyCommit(t, "feat: this is a new feature")

	tagCmd := newTagCmd(noChangesPushed(), os.Stdout)
	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	require.Equal(t, 4, numHooksExecuted(t))
	assert.FileExists(t, BeforeFile)
	assert.FileExists(t, BeforeTagFile)
	assert.FileExists(t, AfterTagFile)
	assert.FileExists(t, AfterFile)
}
