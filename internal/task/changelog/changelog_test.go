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

package changelog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_NoNextTag(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommits(t, "first commit", "second commit", "third commit")

	err := Task{}.Run(&context.Context{})
	require.NoError(t, err)

	assert.False(t, changelogExists(t))
}

func changelogExists(t *testing.T) bool {
	t.Helper()

	current, err := os.Getwd()
	require.NoError(t, err)

	if _, err := os.Stat(filepath.Join(current, MarkdownFile)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		require.NoError(t, err)
	}

	return true
}

func TestRun_CreatedIfNotExists(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommitsAndTag(t, "1.0.0", "first commit", "second commit")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	assert.True(t, changelogExists(t))
}

func TestRun_Staged(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommitsAndTag(t, "1.0.0", "first commit", "second commit")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	stg, _ := git.Staged()
	assert.Len(t, stg, 1)
	assert.Equal(t, "CHANGELOG.md", stg[0])
}

func TestRun_AppendToUnsupportedTemplate(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommitsAndTag(t, "1.0.0", "first commit")

	cl := `# Changelog
This changelog is deliberately missing the append marker`
	ioutil.WriteFile(MarkdownFile, []byte(cl), 0644)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
	}

	err := Task{}.Run(ctx)
	require.ErrorIs(t, err, ErrNoAppendHeader)
}

func TestRun_AppendToExisting(t *testing.T) {
	ih := git.InitRepo(t)
	h1 := git.EmptyCommitsAndTag(t, "1.0.0", "first commit")
	h2 := git.EmptyCommitsAndTag(t, "1.1.0", "second commit", "third commit")

	// Initial changelog
	cl := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 1.0.0 - 2021-09-17

- %s first commit
- %s initialise repo
`, changelogHash(t, h1[0]), changelogHash(t, ih))
	ioutil.WriteFile(MarkdownFile, []byte(cl), 0644)

	ctx := &context.Context{
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 1.1.0 - %s

- %s third commit
- %s second commit

## 1.0.0 - 2021-09-17

- %s first commit
- %s %s
`, changelogDate(t), changelogHash(t, h2[1]), changelogHash(t, h2[0]), changelogHash(t, h1[0]),
		changelogHash(t, ih), git.InitCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_EntriesFromFirstTag(t *testing.T) {
	ih := git.InitRepo(t)
	h := git.EmptyCommitsAndTag(t, "1.0.0", "first commit", "second commit")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 1.0.0 - %s

- %s second commit
- %s first commit
- %s %s
`, changelogDate(t), changelogHash(t, h[1]), changelogHash(t, h[0]), changelogHash(t, ih), git.InitCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func readChangelog(t *testing.T) string {
	t.Helper()

	data, err := ioutil.ReadFile(MarkdownFile)
	require.NoError(t, err)

	return string(data)
}

func changelogDate(t *testing.T) string {
	t.Helper()
	return time.Now().UTC().Format(ChangeDate)
}

func changelogHash(t *testing.T, hash string) string {
	t.Helper()
	return fmt.Sprintf("`%s`", hash)
}

func TestRun_DiffOnly(t *testing.T) {
	git.InitRepo(t)
	git.Tag("1.0.0")
	h := git.EmptyCommitsAndTag(t, "1.1.0", "first commit", "second commit", "third commit")

	var buf bytes.Buffer
	ctx := &context.Context{
		Out:           &buf,
		ChangelogDiff: true,
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`## 1.1.0 - %s

- %s third commit
- %s second commit
- %s first commit
`, changelogDate(t), changelogHash(t, h[2]), changelogHash(t, h[1]), changelogHash(t, h[0]))

	assert.False(t, changelogExists(t))
	assert.Equal(t, expected, buf.String())
}

func TestRun_NoLogEntries(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommitAndTag(t, "1.0.0", "commit")

	err := git.Tag("2.0.0")
	require.NoError(t, err)

	ctx := &context.Context{
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "2.0.0",
		},
	}

	err = Task{}.Run(ctx)
	require.NoError(t, err)
	assert.False(t, changelogExists(t))
}

func TestRun_WithExcludes(t *testing.T) {
	git.InitRepo(t)
	git.Tag("1.0.0")
	h := git.EmptyCommitsAndTag(t, "1.1.0", "first commit", "exclude: second commit", "third commit", "ignore: forth commit")

	var buf bytes.Buffer
	ctx := &context.Context{
		Out:               &buf,
		ChangelogDiff:     true,
		ChangelogExcludes: []string{"exclude", "ignore"},
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`## 1.1.0 - %s

- %s third commit
- %s first commit
`, changelogDate(t), changelogHash(t, h[2]), changelogHash(t, h[0]))

	assert.False(t, changelogExists(t))
	assert.Equal(t, expected, buf.String())
}

func TestRun_ExcludeAllEntries(t *testing.T) {
	git.InitRepo(t)
	git.Tag("1.0.0")
	git.EmptyCommitsAndTag(t, "1.1.0", "prefix: first commit", "prefix: second commit", "prefix: third commit", "prefix: forth commit")

	var buf bytes.Buffer
	ctx := &context.Context{
		Out:               &buf,
		ChangelogExcludes: []string{"prefix"},
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	assert.False(t, changelogExists(t))
	assert.Equal(t, "", buf.String())
}

func TestRun_AllTags(t *testing.T) {
	ih := git.InitRepo(t)
	h1 := git.EmptyCommitAndTag(t, "0.1.0", "feat: first feature")
	h2 := git.EmptyCommitAndTag(t, "0.2.0", "fix: first bug")
	h3 := git.EmptyCommitAndTag(t, "0.3.0", "refactor: use go embed")

	ctx := &context.Context{
		ChangelogAll: true,
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 0.3.0 - %[1]s

- %[2]s refactor: use go embed

## 0.2.0 - %[1]s

- %[3]s fix: first bug

## 0.1.0 - %[1]s

- %[4]s feat: first feature
- %[5]s %[6]s
`, changelogDate(t), changelogHash(t, h3), changelogHash(t, h2), changelogHash(t, h1), changelogHash(t, ih), git.InitCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_AllTagsDiffOnly(t *testing.T) {
	ih := git.InitRepo(t)
	h1 := git.EmptyCommitAndTag(t, "0.1.0", "feat: first feature")
	h2 := git.EmptyCommitAndTag(t, "0.2.0", "fix: first bug")
	h3 := git.EmptyCommitAndTag(t, "0.3.0", "refactor: use go embed")

	var buf bytes.Buffer
	ctx := &context.Context{
		ChangelogAll:  true,
		ChangelogDiff: true,
		Out:           &buf,
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`## 0.3.0 - %[1]s

- %[2]s refactor: use go embed

## 0.2.0 - %[1]s

- %[3]s fix: first bug

## 0.1.0 - %[1]s

- %[4]s feat: first feature
- %[5]s %[6]s
`, changelogDate(t), changelogHash(t, h3), changelogHash(t, h2), changelogHash(t, h1), changelogHash(t, ih), git.InitCommit)

	assert.False(t, changelogExists(t))
	assert.Equal(t, expected, buf.String())
}

func TestRun_AllWithExcludes(t *testing.T) {
	ih := git.InitRepo(t)
	h1 := git.EmptyCommitAndTag(t, "0.1.0", "feat: first feature")
	git.EmptyCommitAndTag(t, "0.2.0", "fix: first bug")
	h2 := git.EmptyCommitAndTag(t, "0.3.0", "refactor: use go embed")

	ctx := &context.Context{
		ChangelogAll:      true,
		ChangelogExcludes: []string{"fix"},
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 0.3.0 - %[1]s

- %[2]s refactor: use go embed

## 0.2.0 - %[1]s

## 0.1.0 - %[1]s

- %[3]s feat: first feature
- %[4]s %[5]s
`, changelogDate(t), changelogHash(t, h2), changelogHash(t, h1), changelogHash(t, ih), git.InitCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_SortCommitsAscending(t *testing.T) {
	ih := git.InitRepo(t)
	h1 := git.EmptyCommitsAndTag(t, "1.0.0", "docs: update to docs", "fix: first bug", "feat: first feature")
	h2 := git.EmptyCommitsAndTag(t, "2.0.0", "ci: tweak scripts", "feat: next feature")

	ctx := &context.Context{
		ChangelogAll:  true,
		ChangelogSort: "asc",
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 2.0.0 - %[1]s

- %[2]s ci: tweak scripts
- %[3]s feat: next feature

## 1.0.0 - %[1]s

- %[4]s %[5]s
- %[6]s docs: update to docs
- %[7]s fix: first bug
- %[8]s feat: first feature
`, changelogDate(t),
		changelogHash(t, h2[0]),
		changelogHash(t, h2[1]),
		changelogHash(t, ih),
		git.InitCommit,
		changelogHash(t, h1[0]),
		changelogHash(t, h1[1]),
		changelogHash(t, h1[2]))

	assert.Equal(t, expected, readChangelog(t))
}
