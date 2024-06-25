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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangelog(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLog("(tag: 0.1.0) feature"))

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))
}

func TestChangelog_DiffOnly(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLog("(tag: 0.1.0) feature"))

	var buf bytes.Buffer

	chglogCmd := newChangelogCmd(noChangesPushed(), &buf)
	chglogCmd.Cmd.SetArgs([]string{"--diff-only"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.False(t, changelogExists(t))
	assert.Contains(t, buf.String(), "## 0.1.0")
}

func TestChangelog_WithExclude(t *testing.T) {
	log := `(tag: 2.0.0) docs: some new docs
ci: a ci task
fix: a new fix
feat: a new feat`
	gittest.InitRepository(t, gittest.WithLog(log))

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--exclude", "^ci,^docs"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, "feat: a new feat")
	assert.Contains(t, cl, "fix: a new fix")
	assert.NotContains(t, cl, "ci: a ci task")
	assert.NotContains(t, cl, "docs: some new docs")
}

func TestChangelog_WithInclude(t *testing.T) {
	log := `(tag: 2.0.0) docs: some new docs
ci: a ci task
fix(scope): a new fix
feat(scope): a new feat`
	gittest.InitRepository(t, gittest.WithLog(log))

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--include", "^.*\\(scope\\)"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, "feat(scope): a new feat")
	assert.Contains(t, cl, "fix(scope): a new fix")
	assert.NotContains(t, cl, "ci: a ci task")
	assert.NotContains(t, cl, "docs: some new docs")
}

func TestChangelog_All(t *testing.T) {
	log := `(tag: 0.5.0) feature 5
(tag: 0.4.0) feature 4
(tag: 0.3.0) feature 3
(tag: 0.2.0) feature 2
(tag: 0.1.0) feature 1`
	gittest.InitRepository(t, gittest.WithLog(log))

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
	log := `(tag: 0.5.0) feature 5
(tag: 0.4.0) feature 4
(tag: 0.3.0) feature 3
(tag: 0.2.0) feature 2
(tag: 0.1.0) feature 1`
	gittest.InitRepository(t, gittest.WithLog(log))

	var buf bytes.Buffer

	chglogCmd := newChangelogCmd(noChangesPushed(), &buf)
	chglogCmd.Cmd.SetArgs([]string{"--all", "--diff-only"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.False(t, changelogExists(t))

	cl := buf.String()
	assert.Contains(t, cl, "## 0.1.0")
	assert.Contains(t, cl, "## 0.2.0")
	assert.Contains(t, cl, "## 0.3.0")
	assert.Contains(t, cl, "## 0.4.0")
	assert.Contains(t, cl, "## 0.5.0")
}

func TestChangelog_SortOrder(t *testing.T) {
	tests := []struct {
		name     string
		sort     string
		commits  []string
		expected []string
	}{
		{
			name:     "Ascending",
			sort:     "asc",
			commits:  []string{"feat: one", "feat: two", "feat: three"},
			expected: []string{"feat: one", "feat: two", "feat: three"},
		},
		{
			name:     "AscendingUpper",
			sort:     "ASC",
			commits:  []string{"feat: one", "feat: two", "feat: three"},
			expected: []string{"feat: one", "feat: two", "feat: three"},
		},
		{
			name:     "Descending",
			sort:     "desc",
			commits:  []string{"feat: one", "feat: two", "feat: three"},
			expected: []string{"feat: three", "feat: two", "feat: one"},
		},
		{
			name:     "DescendingUpper",
			sort:     "DESC",
			commits:  []string{"feat: one", "feat: two", "feat: three"},
			expected: []string{"feat: three", "feat: two", "feat: one"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gittest.InitRepository(t)
			for _, commit := range tt.commits {
				gittest.CommitEmpty(t, commit)
			}
			gittest.Tag(t, "1.0.0")

			chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
			chglogCmd.Cmd.SetArgs([]string{"--sort", tt.sort})

			err := chglogCmd.Cmd.Execute()
			require.NoError(t, err)

			assert.True(t, changelogExists(t))
			cl := readChangelog(t)

			regx := buildChangelogRegex(t, tt.expected)
			assert.True(t, regx.MatchString(cl))
		})
	}
}

func buildChangelogRegex(t *testing.T, commits []string) *regexp.Regexp {
	m := strings.Builder{}
	m.WriteString(fmt.Sprintf("(?im).*%s\n", commits[0]))

	if len(commits) > 1 {
		for i := 1; i < len(commits); i++ {
			m.WriteString(fmt.Sprintf(".*%s\n", commits[i]))
		}
	}

	regx, err := regexp.Compile(m.String())
	require.NoError(t, err)

	return regx
}

func TestChangelog_ExcludesUpliftCommitByDefault(t *testing.T) {
	log := `(tag: 0.1.0) ci(uplift): uplifted version 0.1.0
feat: a new feature
ci: tweak workflow`
	gittest.InitRepository(t, gittest.WithLog(log))

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.NotContains(t, cl, "ci(uplift): uplifted version 0.1.0")
	assert.Contains(t, cl, "feat: a new feature")
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
	gittest.InitRepository(t, gittest.WithLog("feat: this is a new feature"))
	configWithHooks(t)

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	require.Equal(t, 4, numHooksExecuted(t))
	assert.FileExists(t, BeforeFile)
	assert.FileExists(t, BeforeChangelogFile)
	assert.FileExists(t, AfterChangelogFile)
	assert.FileExists(t, AfterFile)
}

func TestChangelog_WithMultiline(t *testing.T) {
	log := `> (tag: 2.0.0) feat: this is a multiline commit
The entire contents of this commit should exist in the changelog.

Multiline formatting should be correct for rendering in markdown
> fix: this is a bug fix
> docs: update documentation
this now includes code examples`
	gittest.InitRepository(t, gittest.WithLog(log))

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--multiline"})

	err := chglogCmd.Cmd.Execute()
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

func TestChangelog_SkipPrerelease(t *testing.T) {
	log := `(tag: 0.1.0) feat: 3
fix: 1
(tag: 0.1.0-pre.2) feat: 2
(tag: 0.1.0-pre.1) feat: 1`
	gittest.InitRepository(t, gittest.WithLog(log))

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--skip-prerelease"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, "## 0.1.0")
	assert.NotContains(t, cl, "## 0.1.0-pre.2")
	assert.NotContains(t, cl, "## 0.1.0-pre.1")
}

func TestChangelog_TrimHeader(t *testing.T) {
	log := `>(tag: 0.1.0) feat: this is a commit
>this line that should be ignored
this line that should also be ignored
feat: second commit`
	gittest.InitRepository(t, gittest.WithLog(log))

	chglogCmd := newChangelogCmd(noChangesPushed(), os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--trim-header"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))

	cl := readChangelog(t)
	assert.Contains(t, cl, `feat: this is a commit`)
	assert.Contains(t, cl, "feat: second commit")
	assert.NotContains(t, cl, "this line that should be ignored")
	assert.NotContains(t, cl, "this line that should also be ignored")
}
