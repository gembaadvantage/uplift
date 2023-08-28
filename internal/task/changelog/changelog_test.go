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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/purpleclay/gitz/gittest"
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
		{
			name: "SkipPrerelease",
			ctx: &context.Context{
				Changelog: context.Changelog{
					SkipPrerelease: true,
				},
				NextVersion: semver.Version{
					Prerelease: "pre.1",
				},
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
	log := `third commit
second commit
first commit`
	gittest.InitRepository(t, gittest.WithLog(log))

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
	log := `(tag: 1.0.0) second commit
first commit`
	gittest.InitRepository(t, gittest.WithLog(log))

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
	log := `(tag: 1.0.0) second commit
first commit`
	gittest.InitRepository(t, gittest.WithLog(log))

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	stg := gittest.PorcelainStatus(t)
	assert.Len(t, stg, 1)
	assert.Equal(t, "A  CHANGELOG.md", stg[0])
}

func TestRun_NotStaged(t *testing.T) {
	log := `(tag: 1.0.0) second commit
first commit`
	gittest.InitRepository(t, gittest.WithLog(log))

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
		NoStage: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	stg := gittest.PorcelainStatus(t)
	assert.Len(t, stg, 1)
	assert.Equal(t, "?? CHANGELOG.md", stg[0])
}

func TestRun_AppendToUnsupportedTemplate(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLog("(tag: 1.0.0) first commit"))

	cl := `# Changelog
This changelog is deliberately missing the append marker`
	os.WriteFile(MarkdownFile, []byte(cl), 0o644)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
	}

	err := Task{}.Run(ctx)
	require.ErrorIs(t, err, ErrNoAppendHeader)
}

func TestRun_AppendToExisting(t *testing.T) {
	log := `(tag: 1.1.0) third commit
second commit
(tag: 1.0.0) first commit`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	// Initial changelog
	cl := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 1.0.0 - 2021-09-17

- %s first commit
- %s %s
`, hashes["first commit"], hashes[gittest.InitialCommit], gittest.InitialCommit)
	os.WriteFile(MarkdownFile, []byte(cl), 0o644)

	ctx := &context.Context{
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
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
`, changelogDate(t), hashes["third commit"], hashes["second commit"], hashes["first commit"],
		hashes[gittest.InitialCommit], gittest.InitialCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func hashLookup(t *testing.T, log []gittest.LogEntry) map[string]string {
	t.Helper()

	hashes := map[string]string{}
	for _, l := range log {
		hashes[l.Message] = fmt.Sprintf("`%s`", l.AbbrevHash)
	}
	return hashes
}

func TestRun_EntriesFromFirstTag(t *testing.T) {
	log := `(tag: 1.1.0) second commit
first commit`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 1.1.0 - %s

- %s second commit
- %s first commit
- %s %s
`, changelogDate(t), hashes["second commit"], hashes["first commit"], hashes[gittest.InitialCommit], gittest.InitialCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func readChangelog(t *testing.T) string {
	t.Helper()

	data, err := os.ReadFile(MarkdownFile)
	require.NoError(t, err)

	return string(data)
}

func changelogDate(t *testing.T) string {
	t.Helper()
	return time.Now().UTC().Format(ChangeDate)
}

func TestRun_DiffOnly(t *testing.T) {
	log := `(tag: 1.1.0) third commit
second commit
first commit
(tag: 1.0.0) won't appear in changelog`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

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
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`## 1.1.0 - %s

- %s third commit
- %s second commit
- %s first commit
`, changelogDate(t), hashes["third commit"], hashes["second commit"], hashes["first commit"])

	assert.False(t, changelogExists(t))
	assert.Equal(t, expected, buf.String())
}

func TestRun_NoLogEntries(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLog("(tag: 1.0.0, tag: 2.0.0) commit"))

	ctx := &context.Context{
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "2.0.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)
	assert.False(t, changelogExists(t))
}

func TestRun_WithExcludes(t *testing.T) {
	log := `(tag: 1.1.0) ignore: forth commit
third commit
exclude(scope): second commit
first commit
(tag: 1.0.0) won't appear in changelog`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		Changelog: context.Changelog{
			DiffOnly: true,
			Exclude:  []string{`^exclude\(scope\)`, "ignore:"},
		},
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`## 1.1.0 - %s

- %s third commit
- %s first commit
`, changelogDate(t), hashes["third commit"], hashes["first commit"])

	assert.False(t, changelogExists(t))
	assert.Equal(t, expected, buf.String())
}

func TestRun_AllTags(t *testing.T) {
	log := `(tag: 0.3.0) third feature
(tag: 0.2.0) second feature
(tag: 0.1.0) first feature`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	ctx := &context.Context{
		Changelog: context.Changelog{
			All: true,
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 0.3.0 - %s

- %s third feature

## 0.2.0 - %s

- %s second feature

## 0.1.0 - %s

- %s first feature
- %s %s
`, changelogDate(t), hashes["third feature"], changelogDate(t), hashes["second feature"], changelogDate(t),
		hashes["first feature"], hashes[gittest.InitialCommit], gittest.InitialCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_AllTagsDiffOnly(t *testing.T) {
	log := `(tag: 0.3.0) third feature
(tag: 0.2.0) second feature
(tag: 0.1.0) first feature`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	var buf bytes.Buffer
	ctx := &context.Context{
		Changelog: context.Changelog{
			All:      true,
			DiffOnly: true,
		},
		Out: &buf,
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`## 0.3.0 - %s

- %s third feature

## 0.2.0 - %s

- %s second feature

## 0.1.0 - %s

- %s first feature
- %s %s
`, changelogDate(t), hashes["third feature"], changelogDate(t), hashes["second feature"], changelogDate(t),
		hashes["first feature"], hashes[gittest.InitialCommit], gittest.InitialCommit)

	assert.False(t, changelogExists(t))
	assert.Equal(t, expected, buf.String())
}

func TestRun_ExcludeAllEntries(t *testing.T) {
	log := `(tag: 1.1.0) prefix: forth commit
prefix: third commit
prefix: second commit
prefix: first commit
(tag: 1.0.0) commit`
	gittest.InitRepository(t, gittest.WithLog(log))

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

	assert.Equal(t, "", buf.String())
}

func TestRun_AllWithExcludes(t *testing.T) {
	log := `(tag: 0.3.0) refactor: use go embed
(tag: 0.2.0) feat: second feature
(tag: 0.1.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	ctx := &context.Context{
		Changelog: context.Changelog{
			All:     true,
			Exclude: []string{"^refactor:"},
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
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

- %s feat: second feature

## 0.1.0 - %s

- %s feat: first feature
- %s %s
`, changelogDate(t), changelogDate(t), hashes["feat: second feature"], changelogDate(t), hashes["feat: first feature"],
		hashes[gittest.InitialCommit], gittest.InitialCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_SortCommitsAscending(t *testing.T) {
	log := `(tag: 2.0.0) feat: second feature
fix: first bug
docs: update to docs
(tag: 1.0.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	ctx := &context.Context{
		Changelog: context.Changelog{
			All:  true,
			Sort: "asc",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
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
- %s feat: second feature

## 1.0.0 - %s

- %s %s
- %s feat: first feature
`, changelogDate(t), hashes["docs: update to docs"], hashes["fix: first bug"], hashes["feat: second feature"],
		changelogDate(t), hashes[gittest.InitialCommit], gittest.InitialCommit, hashes["feat: first feature"])

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_IdentifiedSCM(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLog("(tag: 0.1.0) feat: first feature"))
	log := gittest.Log(t)

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
			Provider:  context.GitHub,
			TagURL:    "https://test.com/tag/{{.Ref}}",
			CommitURL: "https://test.com/commit/{{.Hash}}",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	tag := "[0.1.0](https://test.com/tag/0.1.0)"
	hash := fmt.Sprintf("[`%s`](https://test.com/commit/%s)", log[0].AbbrevHash, log[0].Hash)
	expected := fmt.Sprintf(`## %s - %s

- %s feat: first feature`, tag, changelogDate(t), hash)

	assert.Contains(t, buf.String(), expected)
}

func TestRun_WithMultipleIncludes(t *testing.T) {
	log := `(tag: 1.1.0) fix(common): another fix
feat(scope1): a feature
fix(scope1): a fix
ci: tweak
(tag: 1.0.0) not included in changelog`
	gittest.InitRepository(t, gittest.WithLog(log))

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		Changelog: context.Changelog{
			DiffOnly: true,
			Include:  []string{`^.*\(scope1\)`, `^.*\(common\)`},
		},
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := buf.String()
	assert.Contains(t, actual, "fix(scope1): a fix")
	assert.Contains(t, actual, "feat(scope1): a feature")
	assert.Contains(t, actual, "fix(common): another fix")
	assert.NotContains(t, actual, "ci: tweak")
}

func TestRun_AllWithIncludes(t *testing.T) {
	log := `(tag: 0.3.0) docs: update docs
ci: tweak
feat: another feature
(tag: 0.2.0) feat: second feature
(tag: 0.1.0) feat: first feature`
	gittest.InitRepository(t, gittest.WithLog(log))

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		Changelog: context.Changelog{
			All:      true,
			DiffOnly: true,
			Include:  []string{"^feat:"},
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := buf.String()
	assert.Contains(t, actual, "feat: another feature")
	assert.Contains(t, actual, "feat: second feature")
	assert.Contains(t, actual, "feat: first feature")
	assert.NotContains(t, actual, "ci: tweak")
	assert.NotContains(t, actual, "docs: update docs")
}

func TestRun_CombinedIncludeAndExclude(t *testing.T) {
	log := `(tag: 1.1.0) feat(scope2): another feature
feat(scope1): a feature
fix(scope1): a fix
ci: tweak
(tag: 1.0.0) not included in changelog`
	gittest.InitRepository(t, gittest.WithLog(log))

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		Changelog: context.Changelog{
			DiffOnly: true,
			Include:  []string{`^.*\(scope1\)`, `^.*\(scope2\)`},
			Exclude:  []string{`^fix`},
		},
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := buf.String()
	assert.Contains(t, actual, "feat(scope1): a feature")
	assert.Contains(t, actual, "feat(scope2): another feature")
	assert.NotContains(t, actual, "ci: tweak")
	assert.NotContains(t, actual, "fix(scope1): a fix")
}

func TestRun_MultilineMessages(t *testing.T) {
	log := `> (tag: 1.1.0) feat: this is a multiline commmit

That should be displayed across multiple lines within the changelog.
It should be formatted as expected.

With the correct indentation for rendering in markdown
> feat: this is a single line commit that remains unchanged
> (tag: 1.0.0) not included in changelog`
	gittest.InitRepository(t, gittest.WithLog(log))
	glog := gittest.Log(t)

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		Changelog: context.Changelog{
			DiffOnly:  true,
			Multiline: true,
		},
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`## 1.1.0 - %s

- %s feat: this is a multiline commmit

  That should be displayed across multiple lines within the changelog.
  It should be formatted as expected.

  With the correct indentation for rendering in markdown
- %s feat: this is a single line commit that remains unchanged
`, changelogDate(t), fmt.Sprintf("`%s`", glog[0].AbbrevHash), fmt.Sprintf("`%s`", glog[1].AbbrevHash))

	assert.Equal(t, expected, buf.String())
}

func TestRun_SkipPrerelease(t *testing.T) {
	log := `(tag: 0.2.0) feat: 3
fix: 2
(tag: 0.2.0-pre.2) feat: 2
(tag: 0.2.0-pre.1) fix: 1
(tag: 0.1.0) feat: 1`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	cl := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 1.0.0 - 2021-09-17

- %s feat: 1
- %s %s
`, hashes["first commit"], hashes[gittest.InitialCommit], gittest.InitialCommit)
	os.WriteFile(MarkdownFile, []byte(cl), 0o644)

	ctx := &context.Context{
		Changelog: context.Changelog{
			SkipPrerelease: true,
		},
		CurrentVersion: semver.Version{
			Prerelease: "pre.2",
			Raw:        "0.1.0-pre.2",
		},
		NextVersion: semver.Version{
			Raw: "0.2.0",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 0.2.0 - %s

- %s feat: 3
- %s fix: 2
- %s feat: 2
- %s fix: 1

## 1.0.0 - 2021-09-17

- %s feat: 1
- %s %s
`, changelogDate(t), hashes["feat: 3"], hashes["fix: 2"], hashes["feat: 2"], hashes["fix: 1"],
		hashes["feat: first feature"], hashes[gittest.InitialCommit], gittest.InitialCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_SkipPrereleaseToHead(t *testing.T) {
	log := `(tag: 0.1.0) feat: 3
fix: 1
(tag: 0.1.0-pre.2) feat: 2
(tag: 0.1.0-pre.1) feat: 1`
	gittest.InitRepository(t, gittest.WithLog(log))
	hashes := hashLookup(t, gittest.Log(t))

	ctx := &context.Context{
		Changelog: context.Changelog{
			SkipPrerelease: true,
		},
		CurrentVersion: semver.Version{
			Prerelease: "pre.2",
			Raw:        "0.1.0-pre.2",
		},
		NextVersion: semver.Version{
			Raw: "0.1.0",
		},
		SCM: context.SCM{
			Provider: context.Unrecognised,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## 0.1.0 - %s

- %s feat: 3
- %s fix: 1
- %s feat: 2
- %s feat: 1
- %s %s
`, changelogDate(t), hashes["feat: 3"], hashes["fix: 1"], hashes["feat: 2"], hashes["feat: 1"],
		hashes[gittest.InitialCommit], gittest.InitialCommit)

	assert.Equal(t, expected, readChangelog(t))
}
