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

func TestString(t *testing.T) {
	assert.Equal(t, "generating changelog", Task{}.String())
}

func TestSkip(t *testing.T) {
	tests := []struct {
		name string
		ctx  *context.Context
	}{
		{
			name: "NoVersionChanged",
			ctx: &context.Context{
				NoVersionChanged: true,
			},
		},
		{
			name: "SkipChangelog",
			ctx: &context.Context{
				SkipChangelog: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, Task{}.Skip(tt.ctx))
		})
	}
}

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
	ioutil.WriteFile(MarkdownFile, []byte(cl), 0o644)

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
`, abbrevHash(t, h1[0]), abbrevHash(t, ih))
	ioutil.WriteFile(MarkdownFile, []byte(cl), 0o644)

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
`, changelogDate(t), abbrevHash(t, h2[1]), abbrevHash(t, h2[0]), abbrevHash(t, h1[0]),
		abbrevHash(t, ih), git.InitCommit)

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
`, changelogDate(t), abbrevHash(t, h[1]), abbrevHash(t, h[0]), abbrevHash(t, ih), git.InitCommit)

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

func abbrevHash(t *testing.T, hash string) string {
	t.Helper()
	return fmt.Sprintf("`%s`", hash[:7])
}

func TestRun_DiffOnly(t *testing.T) {
	git.InitRepo(t)
	git.Tag("1.0.0")
	h := git.EmptyCommitsAndTag(t, "1.1.0", "first commit", "second commit", "third commit")

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		Changelog: context.Changelog{
			DiffOnly: true,
		},
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
`, changelogDate(t), abbrevHash(t, h[2]), abbrevHash(t, h[1]), abbrevHash(t, h[0]))

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
		Out: &buf,
		Changelog: context.Changelog{
			DiffOnly: true,
			Exclude:  []string{"exclude", "ignore"},
		},
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
`, changelogDate(t), abbrevHash(t, h[2]), abbrevHash(t, h[0]))

	assert.False(t, changelogExists(t))
	assert.Equal(t, expected, buf.String())
}

func TestRun_ExcludeAllEntries(t *testing.T) {
	git.InitRepo(t)
	git.Tag("1.0.0")
	git.EmptyCommitsAndTag(t, "1.1.0", "prefix: first commit", "prefix: second commit", "prefix: third commit", "prefix: forth commit")

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		Changelog: context.Changelog{
			Exclude: []string{"prefix"},
		},
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
	th := git.TimeBasedTagSeries(t, []string{"0.1.0", "0.2.0", "0.3.0"})

	ctx := &context.Context{
		Changelog: context.Changelog{
			All: true,
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

## 0.3.0 - %s

- %s feat: 3

## 0.2.0 - %s

- %s feat: 2

## 0.1.0 - %s

- %s feat: 1
- %s %s
`, th[2].CreatorDate, abbrevHash(t, th[2].CommitHash), th[1].CreatorDate, abbrevHash(t, th[1].CommitHash), th[0].CreatorDate,
		abbrevHash(t, th[0].CommitHash), abbrevHash(t, ih), git.InitCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_AllTagsDiffOnly(t *testing.T) {
	ih := git.InitRepo(t)
	th := git.TimeBasedTagSeries(t, []string{"0.1.0", "0.2.0", "0.3.0"})

	var buf bytes.Buffer
	ctx := &context.Context{
		Changelog: context.Changelog{
			All:      true,
			DiffOnly: true,
		},
		Out: &buf,
		SCM: context.SCM{
			Provider: git.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`## 0.3.0 - %s

- %s feat: 3

## 0.2.0 - %s

- %s feat: 2

## 0.1.0 - %s

- %s feat: 1
- %s %s
`, th[2].CreatorDate, abbrevHash(t, th[2].CommitHash), th[1].CreatorDate, abbrevHash(t, th[1].CommitHash), th[0].CreatorDate,
		abbrevHash(t, th[0].CommitHash), abbrevHash(t, ih), git.InitCommit)

	assert.False(t, changelogExists(t))
	assert.Equal(t, expected, buf.String())
}

func TestRun_AllWithExcludes(t *testing.T) {
	ih := git.InitRepo(t)
	th := git.TimeBasedTagSeries(t, []string{"0.1.0", "0.2.0"})

	// Commit that will be excluded
	git.EmptyCommitAndTag(t, "0.3.0", "refactor: use go embed")

	ctx := &context.Context{
		Changelog: context.Changelog{
			All:     true,
			Exclude: []string{"refactor"},
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

## 0.3.0 - %s

## 0.2.0 - %s

- %s feat: 2

## 0.1.0 - %s

- %s feat: 1
- %s %s
`, changelogDate(t), th[1].CreatorDate, abbrevHash(t, th[1].CommitHash), th[0].CreatorDate, abbrevHash(t, th[0].CommitHash),
		abbrevHash(t, ih), git.InitCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_SortCommitsAscending(t *testing.T) {
	ih := git.InitRepo(t)
	th := git.TimeBasedTagSeries(t, []string{"1.0.0"})
	hs := git.EmptyCommitsAndTag(t, "2.0.0", "docs: update to docs", "fix: first bug", "feat: first feature")

	ctx := &context.Context{
		Changelog: context.Changelog{
			All:  true,
			Sort: "asc",
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

## 2.0.0 - %s

- %s docs: update to docs
- %s fix: first bug
- %s feat: first feature

## 1.0.0 - %s

- %s %s
- %s feat: 1
`, changelogDate(t), abbrevHash(t, hs[0]), abbrevHash(t, hs[1]), abbrevHash(t, hs[2]), th[0].CreatorDate, abbrevHash(t, ih),
		git.InitCommit, abbrevHash(t, th[0].CommitHash))

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_IdentifiedSCM(t *testing.T) {
	git.InitRepo(t)
	h := git.EmptyCommitAndTag(t, "0.1.0", "feat: first feature")

	var buf bytes.Buffer
	ctx := &context.Context{
		Changelog: context.Changelog{
			DiffOnly: true,
		},
		Out: &buf,
		NextVersion: semver.Version{
			Raw: "0.1.0",
		},
		SCM: context.SCM{
			Provider:  git.GitHub,
			TagURL:    "https://test.com/tag/{{.Ref}}",
			CommitURL: "https://test.com/commit/{{.Hash}}",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	tag := "[0.1.0](https://test.com/tag/0.1.0)"
	hash := fmt.Sprintf("[`%s`](https://test.com/commit/%s)", h[:7], h)
	expected := fmt.Sprintf(`## %s - %s

- %s feat: first feature`, tag, changelogDate(t), hash)

	assert.Contains(t, buf.String(), expected)
}
