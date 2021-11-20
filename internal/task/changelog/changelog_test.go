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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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

func TestRun_ChangelogCreatedIfNotExists(t *testing.T) {
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

func TestRun_ChangelogStaged(t *testing.T) {
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
