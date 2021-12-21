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

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangelog(t *testing.T) {
	taggedRepo(t)

	cmd := newChangelogCmd(&context.Context{})
	err := cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))
}

func TestChangelog_WriteTagToContext(t *testing.T) {
	tests := []struct {
		name       string
		tags       []string
		currentVer string
		nextVer    string
	}{
		{
			name:       "SingleTag",
			tags:       []string{"1.0.0"},
			currentVer: "",
			nextVer:    "1.0.0",
		},
		{
			name:       "MultipleTags",
			tags:       []string{"1.0.0", "1.1.0", "1.2.0", "1.3.0"},
			currentVer: "1.2.0",
			nextVer:    "1.3.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepoWith(t, tt.tags)
			ctx := &context.Context{}

			cmd := newChangelogCmd(ctx)
			err := cmd.Execute()
			require.NoError(t, err)

			require.Equal(t, tt.currentVer, ctx.CurrentVersion.Raw)
			require.Equal(t, tt.nextVer, ctx.NextVersion.Raw)
		})
	}
}

func TestChangelog_DiffOnly(t *testing.T) {
	taggedRepo(t)

	var buf bytes.Buffer
	ctx := context.New(config.Uplift{}, &buf)

	cmd := newChangelogCmd(ctx)
	cmd.SetArgs([]string{"--diff-only"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.False(t, changelogExists(t))
	assert.NotEmpty(t, buf.String())
	assert.True(t, ctx.ChangelogDiff)
}

func TestChangelog_WithExclude(t *testing.T) {
	taggedRepo(t)

	ctx := context.New(config.Uplift{}, nil)

	cmd := newChangelogCmd(ctx)
	cmd.SetArgs([]string{"--exclude", "prefix1", "--exclude", "prefix2"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))
	assert.Len(t, ctx.ChangelogExcludes, 2)
	assert.Contains(t, ctx.ChangelogExcludes[0], "prefix1")
	assert.Contains(t, ctx.ChangelogExcludes[1], "prefix2")
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
