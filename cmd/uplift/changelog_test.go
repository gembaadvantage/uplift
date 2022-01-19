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
	taggedRepo(t)

	chglogCmd := newChangelogCmd(&globalOptions{}, os.Stdout)
	chglogCmd.Cmd.SetArgs([]string{"--exclude", "prefix1,prefix2"})

	err := chglogCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, changelogExists(t))
	assert.Len(t, chglogCmd.Opts.Exclude, 2)
	assert.Contains(t, chglogCmd.Opts.Exclude[0], "prefix1")
	assert.Contains(t, chglogCmd.Opts.Exclude[1], "prefix2")
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
