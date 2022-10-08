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
	"os"
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
		name     string
		cloneURL string
		host     string
		owner    string
		repo     string
		path     string
	}{
		{
			name:     "GitHubSSH",
			cloneURL: "git@github.com:owner/testing1.git",
			host:     "github.com",
			owner:    "owner",
			repo:     "testing1",
			path:     "owner/testing1",
		},
		{
			name:     "GitHubHTTPS",
			cloneURL: "https://github.com/owner/testing2.git",
			host:     "github.com",
			owner:    "owner",
			repo:     "testing2",
			path:     "owner/testing2",
		},
		{
			name:     "GitHubHTTPSWithAccessToken",
			cloneURL: "https://token@github.com/owner/testing3.git",
			host:     "github.com",
			owner:    "owner",
			repo:     "testing3",
			path:     "owner/testing3",
		},
		{
			name:     "GitLabSSH",
			cloneURL: "git@gitlab.com:owner/testing4.git",
			host:     "gitlab.com",
			owner:    "owner",
			repo:     "testing4",
			path:     "owner/testing4",
		},
		{
			name:     "GitLabHTTPS",
			cloneURL: "https://gitlab.com/owner/testing5.git",
			host:     "gitlab.com",
			owner:    "owner",
			repo:     "testing5",
			path:     "owner/testing5",
		},
		{
			name:     "GitLabHTTPSWithAccessToken",
			cloneURL: "https://oauth:token@gitlab.com/owner/testing6.git",
			host:     "gitlab.com",
			owner:    "owner",
			repo:     "testing6",
			path:     "owner/testing6",
		},
		{
			name:     "GitLabUsernamePasswordHTTPS",
			cloneURL: "https://username:password@gitlab.com/owner/testing7.git",
			host:     "gitlab.com",
			owner:    "owner",
			repo:     "testing7",
			path:     "owner/testing7",
		},
		{
			name:     "GitLabHTTPSWithNestedSubGroups",
			cloneURL: "https://gitlab.com/owner/nested/subgroups/testing8.git",
			host:     "gitlab.com",
			owner:    "owner",
			repo:     "testing8",
			path:     "owner/nested/subgroups/testing8",
		},
		{
			name:     "CodeCommitSSH",
			cloneURL: "ssh://git-codecommit.eu-west-1.amazonaws.com/v1/repos/testing9",
			host:     "git-codecommit.eu-west-1.amazonaws.com",
			owner:    "",
			repo:     "testing9",
			path:     "testing9",
		},
		{
			name:     "CodeCommitHTTPS",
			cloneURL: "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/testing10",
			host:     "git-codecommit.eu-west-1.amazonaws.com",
			owner:    "",
			repo:     "testing10",
			path:     "testing10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitRepo(t)
			RemoteOrigin(t, tt.cloneURL)

			repo, err := Remote()

			require.NoError(t, err)
			require.Equal(t, tt.cloneURL, repo.Origin)
			require.Equal(t, tt.host, repo.Host)
			require.Equal(t, tt.owner, repo.Owner)
			require.Equal(t, tt.repo, repo.Name)
			require.Equal(t, tt.path, repo.Path)
		})
	}
}

func TestRemote_CodeCommitGRC(t *testing.T) {
	InitRepo(t)
	RemoteOrigin(t, "codecommit::eu-west-1://profile@repository")

	repo, err := Remote()

	require.NoError(t, err)
	require.Equal(t, "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository", repo.Origin)
	require.Equal(t, "git-codecommit.eu-west-1.amazonaws.com", repo.Host)
	require.Equal(t, "", repo.Owner)
	require.Equal(t, "repository", repo.Name)
}

func TestRemote_HTTPOrigin(t *testing.T) {
	InitRepo(t)
	RemoteOrigin(t, "http://example.com/owner/repository.git")

	repo, err := Remote()

	require.NoError(t, err)
	require.Equal(t, "http://example.com/owner/repository.git", repo.Origin)
	require.Equal(t, "example.com", repo.Host)
	require.Equal(t, "owner", repo.Owner)
	require.Equal(t, "repository", repo.Name)
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

func TestAllTags_WithPrerelease(t *testing.T) {
	InitRepo(t)

	t1, t2, t3, t4, t5 := "0.1.0-alpha+0001", "0.1.0-beta+0001", "0.1.0-beta+0002", "0.1.0", "0.1.1-beta+0001"
	TimeBasedTagSeries(t, []string{t1, t2, t3, t4, t5})

	tags := AllTags()
	require.Len(t, tags, 5)

	assert.Equal(t, "0.1.1-beta+0001", tags[0].Ref)
	assert.Equal(t, "0.1.0", tags[1].Ref)
	assert.Equal(t, "0.1.0-beta+0002", tags[2].Ref)
	assert.Equal(t, "0.1.0-beta+0001", tags[3].Ref)
	assert.Equal(t, "0.1.0-alpha+0001", tags[4].Ref)
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
		"prod",
	})

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
		"0.1.11",
	}

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
		"v2",
	})

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

func TestLatestTag_WithPrerelease(t *testing.T) {
	InitRepo(t)

	t1, t2, t3, t4, t5 := "0.1.0-alpha+0001", "0.1.0-beta+0001", "0.1.0-beta+0002", "0.1.0", "0.1.1-beta+0001"
	TimeBasedTagSeries(t, []string{t1, t2, t3, t4, t5})

	tag := LatestTag()

	assert.Equal(t, "0.1.1-beta+0001", tag.Ref)
}

func TestDescribeTag(t *testing.T) {
	InitRepo(t)
	Run("tag", "1.0.0")

	desc := DescribeTag("1.0.0")

	assert.Equal(t, "1.0.0", desc.Ref)
	assert.Equal(t, time.Now().Format("2006-01-02"), desc.Created)
}

func TestLog(t *testing.T) {
	InitRepo(t)
	EmptyCommits(t,
		"feat: this is a brand new feature",
		`chore(deps): bump knqyf263/trivy-issue-action from 0.0.3 to 0.0.4

Bumps [knqyf263/trivy-issue-action](https://github.com/knqyf263/trivy-issue-action) from 0.0.3 to 0.0.4.
- [Release notes](https://github.com/knqyf263/trivy-issue-action/releases)
- [Commits](https://github.com/knqyf263/trivy-issue-action/compare/v0.0.3...v0.0.4)`,
		`ci: major change to the github workflow

Some extra detail about the workflow`)

	log, err := Log("")
	require.NoError(t, err)

	// Ensure whitespace matches correctly against git log
	assert.Contains(t, log, "feat: this is a brand new feature")
	assert.Contains(t, log, `chore(deps): bump knqyf263/trivy-issue-action from 0.0.3 to 0.0.4
    
    Bumps [knqyf263/trivy-issue-action](https://github.com/knqyf263/trivy-issue-action) from 0.0.3 to 0.0.4.
    - [Release notes](https://github.com/knqyf263/trivy-issue-action/releases)
    - [Commits](https://github.com/knqyf263/trivy-issue-action/compare/v0.0.3...v0.0.4)`)
	assert.Contains(t, log, `ci: major change to the github workflow
    
    Some extra detail about the workflow`)
}

func TestLog_WithTag(t *testing.T) {
	InitRepo(t)
	EmptyCommitsAndTag(t, "1.0.0", "ci: updated existing ci", "docs: new docs", "feat: first feature")
	EmptyCommit(t, `fix: a new bug fix has been added`)

	log, err := Log("1.0.0")
	require.NoError(t, err)

	assert.NotContains(t, log, "ci: updated existing ci")
	assert.NotContains(t, log, "docs: new docs")
	assert.NotContains(t, log, "feat: first feature")
	assert.Contains(t, log, "fix: a new bug fix has been added")
}

func TestAuthor(t *testing.T) {
	InitRepo(t)
	Run("config", "user.name", "uplift")
	Run("config", "user.email", "uplift@test.com")

	details := Author()
	assert.Equal(t, "uplift", details.Name)
	assert.Equal(t, "uplift@test.com", details.Email)
}

func TestAuthorNoNameSet(t *testing.T) {
	InitRepo(t)
	// Setting it to an empty string is the equivalent of it not existing
	Run("config", "user.name", "")
	Run("config", "user.email", "uplift@test.com")

	details := Author()
	assert.Empty(t, details.Name)
	assert.Equal(t, "uplift@test.com", details.Email)
}

func TestAuthorNoEmailSet(t *testing.T) {
	InitRepo(t)
	// Setting it to an empty string is the equivalent of it not existing
	Run("config", "user.name", "uplift")
	Run("config", "user.email", "")

	details := Author()
	assert.Equal(t, "uplift", details.Name)
	assert.Empty(t, details.Email)
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

func TestConfigSet(t *testing.T) {
	InitRepo(t)

	err := ConfigSet(map[string]string{
		"user.name":  "joe.bloggs",
		"user.email": "joe.bloggs@gmail.com",
	})
	require.NoError(t, err)

	name, _ := Clean(Run("config", "--get", "user.name"))
	email, _ := Clean(Run("config", "--get", "user.email"))

	assert.Equal(t, "joe.bloggs", name)
	assert.Equal(t, "joe.bloggs@gmail.com", email)
}

func TestConfigExists(t *testing.T) {
	InitRepo(t)
	Run("config", "--add", "user.name", "joe.bloggs")

	assert.True(t, ConfigExists("user.name", "joe.bloggs"))
}

func UnstagedFile(t *testing.T) string {
	t.Helper()

	err := os.WriteFile("dummy.txt", []byte("hello, world!"), 0o644)
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

func TestLogBetween_AllMultilineCommits(t *testing.T) {
	c1 := `first commit

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.`
	c2 := `second commit

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.`
	c3 := `third commit

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.`

	InitRepo(t)
	EmptyCommits(t, c1, c2, c3)

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
	os.WriteFile("test1.txt", []byte(`testing`), 0o644)
	Stage("test1.txt")

	os.WriteFile("test2.txt", []byte(`testing`), 0o644)
	Stage("test2.txt")

	stg, err := Staged()
	require.NoError(t, err)

	assert.Len(t, stg, 2)
	assert.ElementsMatch(t, stg, []string{"test1.txt", "test2.txt"})
}

func TestStaged_NoFilesStaged(t *testing.T) {
	InitRepo(t)
	os.WriteFile("test.txt", []byte(`testing`), 0o644)

	stg, err := Staged()
	require.NoError(t, err)

	assert.Len(t, stg, 0)
}
