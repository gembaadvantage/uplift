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
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_AppendToExistingChangelog(t *testing.T) {
	ih := git.InitRepo(t)
	h1 := git.EmptyCommitsAndTag(t, "1.0.0", "first commit")
	h2 := git.EmptyCommitsAndTag(t, "1.1.0", "second commit", "third commit")

	// Initial changelog
	cl := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.0.0] - 2021-09-17

%s first commit
%s initialise repo
`, h1[0], ih)
	ioutil.WriteFile(MarkdownFile, []byte(cl), 0644)

	ctx := &context.Context{
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.1.0] - %s

%s third commit
%s second commit

## [1.0.0] - 2021-09-17

%s first commit
%s %s
`, changelogDate(t), h2[1], h2[0], h1[0], ih, git.InitCommit)

	assert.Equal(t, expected, readChangelog(t))
}

func TestRun_ChangelogEntriesFromFirstTag(t *testing.T) {
	ih := git.InitRepo(t)
	h := git.EmptyCommitsAndTag(t, "1.0.0", "first commit", "second commit")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf(`# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.0.0] - %s

%s second commit
%s first commit
%s %s
`, changelogDate(t), h[1], h[0], ih, git.InitCommit)

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
