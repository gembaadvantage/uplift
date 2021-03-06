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
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsRepo(t *testing.T) {
	InitRepo(t)
	assert.True(t, IsRepo())
}

func TestIsRepoDetectsNonGitRepo(t *testing.T) {
	MkTmpDir(t)
	assert.False(t, IsRepo())
}

func TestLatestTag(t *testing.T) {
	InitRepo(t)

	v1 := "v1.0.0"
	Run("tag", v1)
	v2 := "v2.0.0"
	EmptyCommitAndTag(t, v2, "more work")

	tag := LatestTag()
	assert.Equal(t, v2, tag)
}

func TestLatestTagNoTagsExist(t *testing.T) {
	MkTmpDir(t)

	tag := LatestTag()
	assert.Equal(t, "", tag)
}

func TestLatestTagNoSemanticTags(t *testing.T) {
	InitRepo(t)

	v1 := "v1"
	Run("tag", v1)
	v2 := "latest"
	EmptyCommitAndTag(t, v2, "more work")

	tag := LatestTag()
	assert.Equal(t, "", tag)
}

func TestLatestCommit(t *testing.T) {
	InitRepo(t)

	m := "first commit"
	EmptyCommit(t, m)

	c, err := LatestCommit()
	require.NoError(t, err)

	assert.Equal(t, c.Author, "uplift")
	assert.Equal(t, c.Email, "uplift@test.com")
	assert.Equal(t, c.Message, m)
}

func TestLatestCommitMultipleCommits(t *testing.T) {
	InitRepo(t)

	m := "third commit"
	EmptyCommits(t, "first commit", "second commit", m)

	c, err := LatestCommit()
	require.NoError(t, err)

	assert.Equal(t, c.Message, m)
}

func TestTag(t *testing.T) {
	InitRepo(t)

	v := "v1.0.0"
	err := Tag(v)
	require.NoError(t, err)

	_, err = Run("rev-parse", v)
	require.NoError(t, err)
}

func TestAnnotatedTag(t *testing.T) {
	InitRepo(t)

	v := "v1.0.0"
	cd := CommitDetails{
		Author:  "joe.bloggs",
		Email:   "joe.bloggs@gmail.com",
		Message: "a tag commit message",
	}

	err := AnnotatedTag(v, cd)
	require.NoError(t, err)

	out, _ := Clean(Run("for-each-ref", fmt.Sprintf("refs/tags/%s", v),
		"--format='%(taggername):%(taggeremail):%(contents)'"))

	require.Contains(t, out, fmt.Sprintf("%s:<%s>:%s", cd.Author, cd.Email, cd.Message))
}

func TestStage(t *testing.T) {
	InitRepo(t)
	file := UnstagedFile(t)

	err := Stage(file)
	require.NoError(t, err)

	out, err := Clean(Run("status", "--porcelain", "-uno"))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("A  %s", file), out)
}

func TestCommit(t *testing.T) {
	InitRepo(t)
	StagedFile(t)

	cmt := CommitDetails{
		Author:  "joe.bloggs",
		Email:   "joe.bloggs@gmail.com",
		Message: "first commit",
	}

	err := Commit(cmt)
	require.NoError(t, err)

	out, err := Clean(Run("log", "-1", `--pretty=format:'%an:%ae:%B'`))
	require.NoError(t, err)
	assert.Equal(t, "joe.bloggs:joe.bloggs@gmail.com:first commit", out)
}

func UnstagedFile(t *testing.T) string {
	t.Helper()

	err := ioutil.WriteFile("dummy.txt", []byte("hello, world!"), 0644)
	require.NoError(t, err)

	out, err := Clean(Run("status", "--porcelain"))
	require.NoError(t, err)
	require.Equal(t, "?? dummy.txt", out)

	return "dummy.txt"
}

func StagedFile(t *testing.T) string {
	t.Helper()

	file := UnstagedFile(t)

	_, err := Run("add", file)
	require.NoError(t, err)

	return file
}
