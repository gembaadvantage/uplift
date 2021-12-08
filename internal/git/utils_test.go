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

func TestIsRepo_DetectsNonGitRepo(t *testing.T) {
	MkTmpDir(t)
	assert.False(t, IsRepo())
}

func TestAllTags(t *testing.T) {
	InitRepo(t)

	v1 := "v1.0.0"
	EmptyCommitAndTag(t, v1, "first commit")
	v2 := "v2.0.0"
	EmptyCommitAndTag(t, v2, "second commit")
	v3 := "v3.0.0"
	EmptyCommitAndTag(t, v3, "third commit")

	tags := AllTags()
	assert.Equal(t, []string{v3, v2, v1}, tags)
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

func TestLatestTag_NoTagsExist(t *testing.T) {
	MkTmpDir(t)

	tag := LatestTag()
	assert.Equal(t, "", tag)
}

func TestLatestTag_NoSemanticTags(t *testing.T) {
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

func TestLatestCommit_MultipleCommits(t *testing.T) {
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

func TestLogBetween_TwoTags(t *testing.T) {
	InitRepo(t)
	EmptyCommitAndTag(t, "1.0.0", "first commit")
	EmptyCommitsAndTag(t, "2.0.0", "second commit", "third commit", "forth commit")

	log, err := LogBetween("2.0.0", "1.0.0")
	require.NoError(t, err)

	require.Len(t, log, 3)
	assert.Equal(t, log[0].Message, "forth commit")
	assert.Equal(t, log[1].Message, "third commit")
	assert.Equal(t, log[2].Message, "second commit")
}

func TestLogBetween_PrereleaseTag(t *testing.T) {
	InitRepo(t)
	EmptyCommitAndTag(t, "0.1.0", "first commit")
	EmptyCommitsAndTag(t, "0.2.0-beta1+12345", "second commit", "third commit")

	log, err := LogBetween("0.2.0-beta1+12345", "0.1.0")
	require.NoError(t, err)

	require.Len(t, log, 2)
	assert.Equal(t, log[0].Message, "third commit")
	assert.Equal(t, log[1].Message, "second commit")
}

func TestLogBetween_TwoHashes(t *testing.T) {
	InitRepo(t)
	h := EmptyCommits(t, "first commit", "second commit", "third commit", "forth commit")

	log, err := LogBetween(h[2], h[1])
	require.NoError(t, err)

	require.Len(t, log, 1)
	assert.Equal(t, log[0].Message, "third commit")
}

func TestLogBetween_FromSpecificTag(t *testing.T) {
	InitRepo(t)
	EmptyCommitsAndTag(t, "1.0.0", "first commit", "second commit")
	EmptyCommit(t, "third commit")

	log, err := LogBetween("1.0.0", "")
	require.NoError(t, err)

	require.Len(t, log, 3)
	assert.Equal(t, log[0].Message, "second commit")
	assert.Equal(t, log[1].Message, "first commit")
	assert.Equal(t, log[2].Message, InitCommit)
}

func TestLogBetween_FromSpecificHash(t *testing.T) {
	InitRepo(t)
	h := EmptyCommits(t, "first commit", "second commit", "third commit", "forth commit")

	log, err := LogBetween(h[2], "")
	require.NoError(t, err)

	require.Len(t, log, 4)
	assert.Equal(t, log[0].Message, "third commit")
	assert.Equal(t, log[1].Message, "second commit")
	assert.Equal(t, log[2].Message, "first commit")
	assert.Equal(t, log[3].Message, InitCommit)
}

func TestLogBetween_ToSpecificHash(t *testing.T) {
	InitRepo(t)
	h := EmptyCommits(t, "first commit", "second commit", "third commit", "forth commit")

	log, err := LogBetween("", h[2])
	require.NoError(t, err)

	require.Len(t, log, 1)
	assert.Equal(t, log[0].Message, "forth commit")
}

func TestLogBetween_ToSpecificTag(t *testing.T) {
	InitRepo(t)
	EmptyCommitsAndTag(t, "1.0.0", "first commit", "second commit")
	EmptyCommit(t, "third commit")

	log, err := LogBetween("", "1.0.0")
	require.NoError(t, err)

	require.Len(t, log, 1)
	assert.Equal(t, log[0].Message, "third commit")
}

func TestLogBetween_All(t *testing.T) {
	InitRepo(t)
	EmptyCommits(t, "first commit", "second commit", "third commit")

	log, err := LogBetween("", "")
	require.NoError(t, err)

	require.Len(t, log, 4)
	assert.Equal(t, log[0].Message, "third commit")
	assert.Equal(t, log[1].Message, "second commit")
	assert.Equal(t, log[2].Message, "first commit")
	assert.Equal(t, log[3].Message, InitCommit)
}

func TestLogBetween_ErrorInvalidRevision(t *testing.T) {
	InitRepo(t)

	_, err := LogBetween("1234567", "")
	require.Error(t, err)
}

func TestLogBetween_TwoTagsAtSameCommit(t *testing.T) {
	InitRepo(t)
	EmptyCommitAndTag(t, "1.0.0", "first commit")

	err := Tag("1.1.0")
	require.NoError(t, err)

	log, err := LogBetween("1.1.0", "1.0.0")
	require.NoError(t, err)

	assert.Len(t, log, 0)
}

func TestStaged(t *testing.T) {
	InitRepo(t)
	ioutil.WriteFile("test1.txt", []byte(`testing`), 0644)
	Stage("test1.txt")

	ioutil.WriteFile("test2.txt", []byte(`testing`), 0644)
	Stage("test2.txt")

	stg, err := Staged()
	require.NoError(t, err)

	assert.Len(t, stg, 2)
	assert.ElementsMatch(t, stg, []string{"test1.txt", "test2.txt"})
}

func TestStaged_NoFilesStaged(t *testing.T) {
	InitRepo(t)
	ioutil.WriteFile("test.txt", []byte(`testing`), 0644)

	stg, err := Staged()
	require.NoError(t, err)

	assert.Len(t, stg, 0)
}
