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
	"path/filepath"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangelog(t *testing.T) {
	tests := []struct {
		name      string
		tags      []string
		detectTag string
	}{
		{
			name:      "SingleTag",
			tags:      []string{"1.0.0"},
			detectTag: "## 1.0.0",
		},
		{
			name:      "MultipleTags",
			tags:      []string{"1.0.0", "1.1.0", "1.2.0", "1.3.0"},
			detectTag: "## 1.3.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepoWith(t, tt.tags)

			chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
			err := chglogCmd.Cmd.Execute()
			require.NoError(t, err)

			assert.True(t, changelogExists(t))

			chglog := readChangelog(t)
			assert.Contains(t, chglog, tt.detectTag)
		})
	}
}

func TestChangelog_DiffOnly(t *testing.T) {
	taggedRepo(t, "v0.1.0", "feat: a new feature")

	var buf bytes.Buffer

	chglogCmd := newChangelogCmd(noChangesPushed(), &buf)
	chglogCmd.Cmd.SetArgs([]string{"--diff-only"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.False(t, changelogExists(t))
	assert.Contains(t, buf.String(), "## v0.1.0")
}

func TestChangelog_WithExclude(t *testing.T) {
	taggedRepo(t, "2.0.0", "feat: a new feat", "fix: a new fix", "ci: a ci task", "docs: some new docs")

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--exclude", "ci,docs"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, "feat: a new feat")
	assert.Contains(t, cl, "fix: a new fix")
	assert.NotContains(t, cl, "ci: a ci task")
	assert.NotContains(t, cl, "docs: some new docs")
}

func TestChangelog_All(t *testing.T) {
	tagRepoWith(t, []string{"0.1.0", "0.2.0", "0.3.0", "0.4.0", "0.5.0"})

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--all"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, "## 0.1.0")
	assert.Contains(t, cl, "## 0.2.0")
	assert.Contains(t, cl, "## 0.3.0")
	assert.Contains(t, cl, "## 0.4.0")
	assert.Contains(t, cl, "## 0.5.0")
}

func TestChangelog_AllAsDiff(t *testing.T) {
	tagRepoWith(t, []string{"v0.1.0", "v0.2.0", "v0.3.0", "v0.4.0", "v0.5.0"})

	var buf bytes.Buffer

	chglogCmd := newChangelogCmd(noChangesPushed(), &buf)
	chglogCmd.Cmd.SetArgs([]string{"--all", "--diff-only"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.False(t, changelogExists(t))

	cl := buf.String()
	assert.Contains(t, cl, "## v0.1.0")
	assert.Contains(t, cl, "## v0.2.0")
	assert.Contains(t, cl, "## v0.3.0")
	assert.Contains(t, cl, "## v0.4.0")
	assert.Contains(t, cl, "## v0.5.0")
}

// TODO: test commits are in the right order?

func TestChangelog_SortOrder(t *testing.T) {
	tests := []struct {
		name     string
		sort     string
		expected string
	}{
		{
			name:     "Ascending",
			sort:     "asc",
			expected: "asc",
		},
		{
			name:     "AscendingUpper",
			sort:     "ASC",
			expected: "asc",
		},
		{
			name:     "Descending",
			sort:     "desc",
			expected: "desc",
		},
		{
			name:     "DescendingUpper",
			sort:     "DESC",
			expected: "desc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taggedRepo(t, "1.0.0", "feat: a new feature")

			chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
			chglogCmd.Cmd.SetArgs([]string{"--sort", tt.sort})

			err := chglogCmd.Cmd.Execute()
			require.NoError(t, err)

			assert.True(t, changelogExists(t))
			assert.Equal(t, tt.expected, chglogCmd.Opts.Sort)
		})
	}
}

func TestChangelog_ExcludesUpliftCommitByDefault(t *testing.T) {
	taggedRepo(t, "0.1.0", "ci: tweak workflow", "fix: a bug fix", "ci(uplift): uplifted version 0.1.0")

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.NotContains(t, cl, "ci(uplift): uplifted version 0.1.0")
	assert.Contains(t, cl, "fix: a bug fix")
	assert.Contains(t, cl, "ci: tweak workflow")
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

func TestChangelog_Hooks(t *testing.T) {
	git.InitRepo(t)
	configWithHooks(t)
	git.EmptyCommit(t, "feat: this is a new feature")

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	require.Equal(t, 4, numHooksExecuted(t))
	assert.FileExists(t, BeforeFile)
	assert.FileExists(t, BeforeChangelogFile)
	assert.FileExists(t, AfterChangelogFile)
	assert.FileExists(t, AfterFile)
}
