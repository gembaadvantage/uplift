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

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangelog(t *testing.T) {
	taggedRepo(t)

	chglogCmd := newChangelogCmd(&globalOptions{}, os.Stdout)
	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))
}

func TestChangelog_DetectsTags(t *testing.T) {
	tests := []struct {
		name      string
		tags      []string
		detectTag string
	}{
		{
			name:      "SingleTag",
			tags:      []string{"1.0.0"},
			detectTag: "## [1.0.0]",
		},
		{
			name:      "MultipleTags",
			tags:      []string{"1.0.0", "1.1.0", "1.2.0", "1.3.0"},
			detectTag: "## [1.3.0]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepoWith(t, tt.tags)

			chglogCmd := newChangelogCmd(&globalOptions{}, os.Stdout)
			err := chglogCmd.Cmd.Execute()
			require.NoError(t, err)

			assert.True(t, changelogExists(t))
			chglog := readChangelog(t)
			assert.Contains(t, chglog, tt.detectTag)
		})
	}
}

func TestChangelog_DiffOnly(t *testing.T) {
	taggedRepo(t)

	var buf bytes.Buffer

	chglogCmd := newChangelogCmd(&globalOptions{}, &buf)
	chglogCmd.Cmd.SetArgs([]string{"--diff-only"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.False(t, changelogExists(t))
	assert.NotEmpty(t, buf.String())
	assert.True(t, chglogCmd.Opts.DiffOnly)
}

func TestChangelog_WithExclude(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommitsAndTag(t, "2.0.0",
		"feat: a new feat",
		"fix: a new fix",
		"ci: a ci task",
		"docs: some new docs")

	chglogCmd := newChangelogCmd(&globalOptions{}, os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--exclude", "ci,docs"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))
	assert.Len(t, chglogCmd.Opts.Exclude, 2)
	assert.Contains(t, chglogCmd.Opts.Exclude[0], "ci")
	assert.Contains(t, chglogCmd.Opts.Exclude[1], "docs")

	chglog := readChangelog(t)
	assert.NotContains(t, chglog, "ci:")
	assert.NotContains(t, chglog, "docs:")
}

func TestChangelog_All(t *testing.T) {
	tagRepoWith(t, []string{"0.1.0", "0.2.0", "0.3.0", "0.4.0", "0.5.0"})

	chglogCmd := newChangelogCmd(&globalOptions{}, os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--all"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	chglog := readChangelog(t)
	assert.Contains(t, chglog, "## [0.1.0]")
	assert.Contains(t, chglog, "## [0.2.0]")
	assert.Contains(t, chglog, "## [0.3.0]")
	assert.Contains(t, chglog, "## [0.4.0]")
	assert.Contains(t, chglog, "## [0.5.0]")
}

func TestChangelog_AllAsDiff(t *testing.T) {
	tagRepoWith(t, []string{"0.1.0", "0.2.0", "0.3.0", "0.4.0", "0.5.0"})

	var buf bytes.Buffer

	chglogCmd := newChangelogCmd(&globalOptions{}, &buf)
	chglogCmd.Cmd.SetArgs([]string{"--all", "--diff-only"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.False(t, changelogExists(t))

	chglog := buf.String()
	assert.Contains(t, chglog, "## [0.1.0]")
	assert.Contains(t, chglog, "## [0.2.0]")
	assert.Contains(t, chglog, "## [0.3.0]")
	assert.Contains(t, chglog, "## [0.4.0]")
	assert.Contains(t, chglog, "## [0.5.0]")
}

func changelogExists(t *testing.T) bool {
	t.Helper()

	current, err := os.Getwd()
	require.NoError(t, err)

	if _, err := os.Stat(filepath.Join(current, "CHANGELOG.md")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		require.NoError(t, err)
	}

	return true
}

func readChangelog(t *testing.T) string {
	t.Helper()

	data, err := os.ReadFile("CHANGELOG.md")
	require.NoError(t, err)

	return string(data)
}
