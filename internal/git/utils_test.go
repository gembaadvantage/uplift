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

package git

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsInstalled(t *testing.T) {
	assert.True(t, IsInstalled())
}

func TestIsInstalled_NotOnPath(t *testing.T) {
	// Temporarily blast the PATH variable within this test
	t.Setenv("PATH", "")

	assert.False(t, IsInstalled())
}

func TestIsRepo(t *testing.T) {
	InitRepo(t)
	assert.True(t, IsRepo())
}

func TestIsRepo_DetectsNonGitRepo(t *testing.T) {
	MkTmpDir(t)
	assert.False(t, IsRepo())
}

func TestIsShallow(t *testing.T) {
	InitShallowRepo(t)

	assert.True(t, IsShallow())
}

func TestIsShallow_FullHistory(t *testing.T) {
	InitRepo(t)

	assert.False(t, IsShallow())
}

func TestIsDetached(t *testing.T) {
	InitRepo(t)
	h := EmptyCommit(t, "this is a test")

	_, err := Run("checkout", h)
	require.NoError(t, err)

	assert.True(t, IsDetached())
}

func TestIsDetached_OnBranch(t *testing.T) {
	InitRepo(t)

	assert.False(t, IsDetached())
}

func TestCheckDirty(t *testing.T) {
	InitRepo(t)
	EmptyCommit(t, "this is a test")

	out, err := CheckDirty()
	require.NoError(t, err)
	assert.Empty(t, out)
}

func TestCheckDirty_UnCommitted(t *testing.T) {
	InitRepo(t)

	TouchFiles(t, "main.go", "testing.go")
	Stage("main.go")
	Stage("testing.go")

	out, err := CheckDirty()
	require.NoError(t, err)

	exp := `A  main.go
A  testing.go`
	assert.Equal(t, exp, out)
}

func TestCheckDirty_UnStaged(t *testing.T) {
	InitRepo(t)

	// Add an empty file
	TouchFiles(t, "main.go", "testing.go")

	out, err := CheckDirty()
	require.NoError(t, err)

	exp := `?? main.go
?? testing.go`
	assert.Equal(t, exp, out)
}

func TestRemote(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		scm       SCM
		owner     string
		repo      string
		cloneURL  string
		browseURL string
	}{
		{
			name:      "GitHubSSH",
			url:       "git@github.com:owner/testing1.git",
			scm:       GitHub,
			owner:     "owner",
			repo:      "testing1",
			cloneURL:  "https://github.com/owner/testing1",
			browseURL: "https://github.com/owner/testing1",
		},
		{
			name:      "GitHubHTTPS",
			url:       "https://github.com/owner/testing2.git",
			scm:       GitHub,
			owner:     "owner",
			repo:      "testing2",
			cloneURL:  "https://github.com/owner/testing2",
			browseURL: "https://github.com/owner/testing2",
		},
		{
			name:      "GitHubHTTPSWithAccessToken",
			url:       "https://token@github.com/owner/testing3.git",
			scm:       GitHub,
			owner:     "owner",
			repo:      "testing3",
			cloneURL:  "https://github.com/owner/testing3",
			browseURL: "https://github.com/owner/testing3"},
		{
			name:      "GitLabSSH",
			url:       "git@gitlab.com:owner/testing4.git",
			scm:       GitLab,
			owner:     "owner",
			repo:      "testing4",
			cloneURL:  "https://gitlab.com/owner/testing4",
			browseURL: "https://gitlab.com/owner/testing4",
		},
		{
			name:      "GitLabHTTPS",
			url:       "https://gitlab.com/owner/testing5.git",
			scm:       GitLab,
			owner:     "owner",
			repo:      "testing5",
			cloneURL:  "https://gitlab.com/owner/testing5",
			browseURL: "https://gitlab.com/owner/testing5",
		},
		{
			name:      "GitLabHTTPSWithAccessToken",
			url:       "https://oauth:token@gitlab.com/owner/testing6.git",
			scm:       GitLab,
			owner:     "owner",
			repo:      "testing6",
			cloneURL:  "https://gitlab.com/owner/testing6",
			browseURL: "https://gitlab.com/owner/testing6",
		},
		{
			name:      "GitLabUsernamePasswordHTTPS",
			url:       "https://username:password@gitlab.com/owner/testing7.git",
			scm:       GitLab,
			owner:     "owner",
			repo:      "testing7",
			cloneURL:  "https://gitlab.com/owner/testing7",
			browseURL: "https://gitlab.com/owner/testing7",
		},
		{
			name:      "CodeCommitSSH",
			url:       "ssh://git-codecommit.eu-west-1.amazonaws.com/v1/repos/testing8",
			scm:       CodeCommit,
			owner:     "",
			repo:      "testing8",
			cloneURL:  "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/testing8",
			browseURL: "https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/testing8",
		},
		{
			name:      "CodeCommitHTTPS",
			url:       "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/testing9",
			scm:       CodeCommit,
			owner:     "",
			repo:      "testing9",
			cloneURL:  "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/testing9",
			browseURL: "https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/testing9",
		},
		{
			name:      "CodeCommitGRC",
			url:       "codecommit::eu-west-1://profile@testing10",
			scm:       CodeCommit,
			owner:     "",
			repo:      "testing10",
			cloneURL:  "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/testing10",
			browseURL: "https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/testing10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitRepo(t)
			RemoteOrigin(t, tt.url)

			repo, err := Remote()

			require.NoError(t, err)
			require.Equal(t, tt.scm, repo.Provider)
			require.Equal(t, tt.owner, repo.Owner)
			require.Equal(t, tt.repo, repo.Name)
			require.Equal(t, tt.cloneURL, repo.CloneURL)
			require.Equal(t, tt.browseURL, repo.BrowseURL)
		})
	}
}

func TestRemote_NoRemoteSet(t *testing.T) {
	InitRepo(t)
	Run("remote", "remove", "origin")

	_, err := Remote()
	require.Error(t, err)
}

func TestRemote_MalformedURL(t *testing.T) {
	InitRepo(t)
	RemoteOrigin(t, "whizzbang.com/repository")

	_, err := Remote()
	require.EqualError(t, err, "malformed repository URL: whizzbang.com/repository")
}

func TestRemote_Unrecognised(t *testing.T) {
	InitRepo(t)
	RemoteOrigin(t, "https://whizzbang.com/owner/repository")

	repo, err := Remote()

	require.NoError(t, err)
	assert.Equal(t, Unrecognised, repo.Provider)
}

func TestAllTags(t *testing.T) {
	InitRepo(t)

	v1, v2, v3 := "v1.0.0", "v2.0.0", "v3.0.0"
	TimeBasedTagSeries(t, []string{v1, v2, v3})

	tags := AllTags()
	require.Len(t, tags, 3)
	assert.Equal(t, "v3.0.0", tags[0].Ref)
	assert.Equal(t, "v2.0.0", tags[1].Ref)
	assert.Equal(t, "v1.0.0", tags[2].Ref)
}

func TestAllTags_MixedTagConventions(t *testing.T) {
	InitRepo(t)

	v1, v2, v3 := "1.0.0", "v2.0.0", "3.0.0"
	TimeBasedTagSeries(t, []string{v1, v2, v3})

	tags := AllTags()
	require.Len(t, tags, 3)
	assert.Equal(t, "3.0.0", tags[0].Ref)
	assert.Equal(t, "v2.0.0", tags[1].Ref)
	assert.Equal(t, "1.0.0", tags[2].Ref)
}

func TestAllTags_FiltersNonSemanticTags(t *testing.T) {
	InitRepo(t)

	t1, t2, t3, t4 := "v1", "1.1.0", "in.va.lid", "1.2.0"
	TimeBasedTagSeries(t, []string{t1, t2, t3, t4})

	tags := AllTags()
	require.Len(t, tags, 2)
	assert.Equal(t, "1.2.0", tags[0].Ref)
	assert.Equal(t, "1.1.0", tags[1].Ref)
}

func TestAllTags_CommitWithMultipleTags(t *testing.T) {
	InitRepo(t)

	t1, t2 := "1.0.0", "1.1.0"
	TimeBasedTagSeries(t, []string{t1, t2})
	EmptyCommitAndTags(t, "another commit", "v1", "1.2.0")

	tags := AllTags()
	require.Len(t, tags, 3)
	assert.Equal(t, "1.2.0", tags[0].Ref)
	assert.Equal(t, "1.1.0", tags[1].Ref)
	assert.Equal(t, "1.0.0", tags[2].Ref)
}

func TestAllTags_LargeHistory(t *testing.T) {
	InitRepo(t)

	TimeBasedTagSeries(t, []string{
		"0.1.11",
		"0.1.123",
		"0.9.0",
		"0.10.0",
		"0.12.0",
		"0.123.0",
		"1.0.0",
		"v1",
		"1.1.1",
		"1.1.10",
		"1.10.0",
		"1.11.0",
		"2.0.0",
		"v2",
		"10.1.10",
		"10.11.10",
		"11.0.0",
		"prod"})

	tags := AllTags()
	require.Len(t, tags, 15)

	exp := []string{
		"11.0.0",
		"10.11.10",
		"10.1.10",
		"2.0.0",
		"1.11.0",
		"1.10.0",
		"1.1.10",
		"1.1.1",
		"1.0.0",
		"0.123.0",
		"0.12.0",
		"0.10.0",
		"0.9.0",
		"0.1.123",
		"0.1.11"}

	for i, tag := range tags {
		assert.Equal(t, exp[i], tag.Ref)
	}
}

func TestLatestTag(t *testing.T) {
	InitRepo(t)

	v1, v2 := "v1.0.0", "v2.0.0"
	TimeBasedTagSeries(t, []string{v1, v2})

	tag := LatestTag()
	assert.Equal(t, v2, tag.Ref)
}

func TestLatestTag_LargeHistory(t *testing.T) {
	InitRepo(t)

	TimeBasedTagSeries(t, []string{
		"0.1.0",
		"0.2.0",
		"0.9.0",
		"0.10.0",
		"0.11.0",
		"0.29.0",
		"1.0.0",
		"v1",
		"1.9.0",
		"1.10.0",
		"1.11.1",
		"2.0.0",
		"v2"})

	tag := LatestTag()
	assert.Equal(t, "2.0.0", tag.Ref)
}

func TestLatestTag_NoTagsExist(t *testing.T) {
	MkTmpDir(t)

	tag := LatestTag()
	assert.Equal(t, "", tag.Ref)
}

func TestLatestTag_NoSemanticTags(t *testing.T) {
	InitRepo(t)

	v1 := "v1"
	Run("tag", v1)
	v2 := "latest"
	EmptyCommitAndTag(t, v2, "more work")

	tag := LatestTag()
	assert.Equal(t, "", tag.Ref)
}

func TestLatestTag_MixedTagConventions(t *testing.T) {
	InitRepo(t)

	v1, v2, v3 := "v1.0.0", "2.0.0", "v3.0.0"
	TimeBasedTagSeries(t, []string{v1, v2, v3})

	tag := LatestTag()
	assert.Equal(t, "v3.0.0", tag.Ref)
}

func TestLatestTag_CommitWithMixedTags(t *testing.T) {
	InitRepo(t)

	EmptyCommitAndTags(t, "commit", "v1", "prod", "1.0.0")

	tag := LatestTag()
	assert.Equal(t, "1.0.0", tag.Ref)
}

func TestDescribeTag(t *testing.T) {
	InitRepo(t)
	Run("tag", "1.0.0")

	desc := DescribeTag("1.0.0")

	assert.Equal(t, "1.0.0", desc.Ref)
	assert.Equal(t, time.Now().Format("2006-01-02"), desc.Created)
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

func TestCommitDetails_String(t *testing.T) {
	cd := CommitDetails{
		Author:  "uplift",
		Email:   "uplift@test.com",
		Message: "this is a test commit",
	}

	assert.Equal(t, "uplift <uplift@test.com>\nthis is a test commit", cd.String())
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

	log, err := LogBetween("2.0.0", "1.0.0", []string{})
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

	log, err := LogBetween("0.2.0-beta1+12345", "0.1.0", []string{})
	require.NoError(t, err)

	require.Len(t, log, 2)
	assert.Equal(t, log[0].Message, "third commit")
	assert.Equal(t, log[1].Message, "second commit")
}

func TestLogBetween_TwoHashes(t *testing.T) {
	InitRepo(t)
	h := EmptyCommits(t, "first commit", "second commit", "third commit", "forth commit")

	log, err := LogBetween(h[2], h[1], []string{})
	require.NoError(t, err)

	require.Len(t, log, 1)
	assert.Equal(t, log[0].Message, "third commit")
}

func TestLogBetween_FromSpecificTag(t *testing.T) {
	InitRepo(t)
	EmptyCommitsAndTag(t, "1.0.0", "first commit", "second commit")
	EmptyCommit(t, "third commit")

	log, err := LogBetween("1.0.0", "", []string{})
	require.NoError(t, err)

	require.Len(t, log, 3)
	assert.Equal(t, log[0].Message, "second commit")
	assert.Equal(t, log[1].Message, "first commit")
	assert.Equal(t, log[2].Message, InitCommit)
}

func TestLogBetween_FromSpecificHash(t *testing.T) {
	InitRepo(t)
	h := EmptyCommits(t, "first commit", "second commit", "third commit", "forth commit")

	log, err := LogBetween(h[2], "", []string{})
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

	log, err := LogBetween("", h[2], []string{})
	require.NoError(t, err)

	require.Len(t, log, 1)
	assert.Equal(t, log[0].Message, "forth commit")
}

func TestLogBetween_ToSpecificTag(t *testing.T) {
	InitRepo(t)
	EmptyCommitsAndTag(t, "1.0.0", "first commit", "second commit")
	EmptyCommit(t, "third commit")

	log, err := LogBetween("", "1.0.0", []string{})
	require.NoError(t, err)

	require.Len(t, log, 1)
	assert.Equal(t, log[0].Message, "third commit")
}

func TestLogBetween_All(t *testing.T) {
	InitRepo(t)
	EmptyCommits(t, "first commit", "second commit", "third commit")

	log, err := LogBetween("", "", []string{})
	require.NoError(t, err)

	require.Len(t, log, 4)
	assert.Equal(t, log[0].Message, "third commit")
	assert.Equal(t, log[1].Message, "second commit")
	assert.Equal(t, log[2].Message, "first commit")
	assert.Equal(t, log[3].Message, InitCommit)
}

func TestLogBetween_ErrorInvalidRevision(t *testing.T) {
	InitRepo(t)

	_, err := LogBetween("1234567", "", []string{})
	require.Error(t, err)
}

func TestLogBetween_TwoTagsAtSameCommit(t *testing.T) {
	InitRepo(t)
	EmptyCommitAndTag(t, "1.0.0", "first commit")

	err := Tag("1.1.0")
	require.NoError(t, err)

	log, err := LogBetween("1.1.0", "1.0.0", []string{})
	require.NoError(t, err)

	assert.Len(t, log, 0)
}

func TestLogBetween_Excludes(t *testing.T) {
	InitRepo(t)
	EmptyCommitAndTag(t, "0.1.0", "first commit")
	EmptyCommits(t, "second commit", "exclude: third commit", "ignore: forth commit")
	EmptyCommitAndTag(t, "0.2.0", "fifth commit")

	log, err := LogBetween("0.2.0", "0.1.0", []string{"exclude", "ignore"})
	require.NoError(t, err)

	assert.Len(t, log, 2)
	assert.Equal(t, log[0].Message, "fifth commit")
	assert.Equal(t, log[1].Message, "second commit")
}

func TestLogBetween_ExcludesWildcards(t *testing.T) {
	InitRepo(t)
	EmptyCommitAndTag(t, "0.1.0", "first commit")
	EmptyCommits(t, "second commit", "exclude: third commit", "exclude(scope): forth commit")
	EmptyCommitAndTag(t, "0.2.0", "fifth commit")

	log, err := LogBetween("0.2.0", "0.1.0", []string{"exclude*"})
	require.NoError(t, err)

	assert.Len(t, log, 2)
	assert.Equal(t, log[0].Message, "fifth commit")
	assert.Equal(t, log[1].Message, "second commit")
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
