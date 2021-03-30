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

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Use test suite and automatically run this before every test

func TestLatestTagSingleTag(t *testing.T) {
	InitRepo(t)

	v := "v1.2.3"
	Run("tag", v)

	tag, err := LatestTag()
	require.NoError(t, err)

	assert.Equal(t, v, tag)
}

func TestLatestTag(t *testing.T) {
	InitRepo(t)

	v1 := "v1.0.0"
	Run("tag", v1)
	v2 := "v2.0.0"
	EmptyCommit(t, "more work")
	Run("tag", v2)

	tag, err := LatestTag()
	require.NoError(t, err)

	assert.Equal(t, v2, tag)
}

func TestLatestTagNoTagsExist(t *testing.T) {
	InitRepo(t)

	_, err := LatestTag()
	require.ErrorIs(t, err, ErrNoTag)
}

func TestLatestCommitMessage(t *testing.T) {
	InitRepo(t)

	m := "first commit"
	EmptyCommit(t, m)

	msg, err := LatestCommitMessage()
	require.NoError(t, err)

	assert.Equal(t, m, msg)
}

func TestLatestCommitMessageMultipleCommits(t *testing.T) {
	InitRepo(t)

	EmptyCommit(t, "first commit")
	EmptyCommit(t, "second commit")
	m := "third commit"
	EmptyCommit(t, m)

	msg, err := LatestCommitMessage()
	require.NoError(t, err)

	assert.Equal(t, m, msg)
}

func TestTag(t *testing.T) {
	InitRepo(t)

	v := "v1.0.0"
	_, err := Tag(v)
	require.NoError(t, err)

	_, err = Run("rev-parse", v)
	require.NoError(t, err)
}
