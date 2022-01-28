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
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

// SCM ...
type SCM string

const (
	GitHub       SCM = "GitHub"
	GitLab       SCM = "GitLab"
	CodeCommit   SCM = "CodeCommit"
	Unrecognised SCM = "Unrecognised"
)

// CommitDetails contains mandatory details about a specific git commit
type CommitDetails struct {
	Message string
	Author  string
	Email   string
}

// LogEntry contains details about a specific git log entry
type LogEntry struct {
	Hash       string
	AbbrevHash string
	Message    string
}

// TagEntry contains details about a specific tag
type TagEntry struct {
	Ref     string
	Created string
}

// Repository ...
type Repository struct {
	Provider SCM
	Owner    string
	Name     string
	URL      string
}

// String prints out a user friendly string representation
func (c CommitDetails) String() string {
	return fmt.Sprintf("%s <%s>\n%s", c.Author, c.Email, c.Message)
}

// Run executes a git command and returns its output or errors
func Run(args ...string) (string, error) {
	var cmd = exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

// IsRepo identifies whether the current working directory is a recognised
// git repository
func IsRepo() bool {
	out, err := Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// Remote retrieves details about the remote origin of a repository
func Remote() (Repository, error) {
	remURL, err := Clean(Run("config", "--get", "remote.origin.url"))
	if err != nil {
		return Repository{}, err
	}

	// Strip off any trailing .git suffix
	rem := strings.TrimRight(remURL, ".git")

	if strings.HasPrefix(rem, "git@") {
		// Sanitise any SSH based URL to ensure it is parseable
		rem = strings.TrimPrefix(rem, "git@")
		rem = strings.Replace(rem, ":", "/", 1)
	} else if strings.HasPrefix(rem, "https://") {
		// Sanitise any HTTPS based URL to ensure it is parseable. Handle username@password inclusion
		rem = rem[strings.LastIndex(rem, ":")+1:]
		rem = strings.TrimPrefix(rem, "//")
	}

	// Special use case for CodeCommit as that prefixes SSH URLs with ssh://
	rem = strings.TrimPrefix(rem, "ssh://")

	u, err := url.Parse(rem)
	if err != nil {
		return Repository{}, err
	}

	// Split into parts
	p := strings.Split(u.Path, "/")
	if len(p) < 3 {
		return Repository{}, fmt.Errorf("malformed repository URL: %s", remURL)
	}

	// No concept of an owner with CodeCommit repositories
	owner := p[1]
	if strings.Contains(p[0], "codecommit") {
		owner = ""
	}

	return Repository{
		Provider: detectSCM(p[0]),
		Owner:    owner,
		Name:     p[len(p)-1],
		URL:      fmt.Sprintf("https://%s", u.Path),
	}, nil
}

func detectSCM(host string) SCM {
	switch host {
	case "github.com":
		return GitHub
	case "gitlab.com":
		return GitLab
	}

	// Handle special case CodeCommit URLs
	if strings.Contains(host, "codecommit") {
		return CodeCommit
	}

	return Unrecognised
}

// FetchTags retrieves all tags associated with the remote repository
func FetchTags() error {
	if _, err := Clean(Run("fetch", "--all", "--tags")); err != nil {
		return err
	}

	return nil
}

// AllTags retrieves all tags within the repository from newest to oldest
func AllTags() []TagEntry {
	return retrieveTags(0)
}

func retrieveTags(n int) []TagEntry {
	tags, err := Clean(Run("for-each-ref",
		"refs/tags/v*.*.*",
		"refs/tags/*.*.*",
		"--sort=-v:refname",
		`--format='%(creatordate:short),%(refname:short)'`,
		fmt.Sprintf("--count=%d", n),
	))
	if err != nil {
		return []TagEntry{}
	}

	// If no tags are found, then just return an empty slice
	if tags == "" {
		return []TagEntry{}
	}

	rows := strings.Split(tags, "\n")
	ents := make([]TagEntry, 0, len(rows))
	for _, r := range rows {
		p := strings.Split(r, ",")
		ents = append(ents, TagEntry{
			Ref:     p[1],
			Created: p[0],
		})
	}

	return ents
}

// LatestTag retrieves the latest tag within the repository
func LatestTag() TagEntry {
	tags := retrieveTags(1)
	if len(tags) == 0 {
		return TagEntry{}
	}

	return tags[0]
}

// DescribeTag retrieves details about a specific tag
func DescribeTag(ref string) TagEntry {
	tag, err := Clean(Run("tag", "-l", ref, `--format='%(creatordate:short),%(refname:short)'`))
	if err != nil {
		return TagEntry{}
	}

	if tag == "" {
		return TagEntry{}
	}

	p := strings.Split(tag, ",")
	return TagEntry{
		Ref:     p[1],
		Created: p[0],
	}
}

// LatestCommit retrieves the latest commit within the repository
func LatestCommit() (CommitDetails, error) {
	out, err := Clean(Run("log", "-1", `--pretty=format:'"%an","%ae","%B"'`))
	if err != nil {
		return CommitDetails{}, err
	}

	// Split the formatted string into its component parts
	p := strings.Split(out, ",")

	// Strip quotes from around each part
	author := p[0][1 : len(p[0])-1]
	email := p[1][1 : len(p[1])-1]
	msg := p[2][1 : len(p[2])-1]

	return CommitDetails{
		Author: author,
		Email:  email,
		// Strip trailing newline
		Message: strings.TrimRight(msg, "\n"),
	}, nil
}

// Tag will create a lightweight tag against the repository and push it to the origin
func Tag(tag string) error {
	if _, err := Clean(Run("tag", tag)); err != nil {
		return err
	}

	// Inspect the repo for an origin. If no origin exists, then skip the push
	if _, err := Clean(Run("remote", "show", "origin")); err != nil {
		return nil
	}

	if _, err := Clean(Run("push", "origin", tag)); err != nil {
		return err
	}

	return nil
}

// AnnotatedTag will create an annotated tag against the repository and push it to the origin
func AnnotatedTag(tag string, cd CommitDetails) error {
	args := []string{
		"-c",
		fmt.Sprintf("user.name='%s'", cd.Author),
		"-c",
		fmt.Sprintf("user.email='%s'", cd.Email),
		"tag",
		"-a",
		tag,
		"-m",
		cd.Message,
	}

	if _, err := Clean(Run(args...)); err != nil {
		return err
	}

	// Inspect the repo for an origin. If no origin exists, then skip the push
	if _, err := Clean(Run("remote", "show", "origin")); err != nil {
		return nil
	}

	if _, err := Clean(Run("push", "origin", tag)); err != nil {
		return err
	}

	return nil
}

// Commit will generate a commit against the repository and push it to the origin.
// The commit will be associated with the provided author and email address
func Commit(cd CommitDetails) error {
	args := []string{
		"-c",
		fmt.Sprintf("user.name='%s'", cd.Author),
		"-c",
		fmt.Sprintf("user.email='%s'", cd.Email),
		"commit",
		"-m",
		cd.Message,
	}

	if _, err := Clean(Run(args...)); err != nil {
		return err
	}

	return nil
}

// Stage will ensure the specified file is staged for the next commit
func Stage(path string) error {
	if _, err := Clean(Run("add", path)); err != nil {
		return err
	}

	return nil
}

// Push all committed changes to the configured origin
func Push() error {
	// Inspect the repo for an origin. If no origin exists, then skip the push
	if _, err := Clean(Run("remote", "show", "origin")); err != nil {
		return nil
	}

	branch, err := Clean(Run("rev-parse", "--abbrev-ref", "HEAD"))
	if err != nil {
		return err
	}

	if _, err := Clean(Run("push", "origin", branch)); err != nil {
		return err
	}

	return nil
}

// LogBetween retrieves all log entries between two points of time within the
// git history of the repository. Supports tags and specific git hashes as its
// reference points. From must always be the closest point to HEAD
func LogBetween(from, to string, excludes []string) ([]LogEntry, error) {
	fmtFrom := from
	if fmtFrom == "" {
		fmtFrom = "HEAD"
	}

	fmtTo := to
	if fmtTo != "" {
		// A range query requires ... ellipses
		fmtTo = fmt.Sprintf("...%s", fmtTo)
	}

	args := []string{
		"log",
		fmt.Sprintf("%s%s", fmtFrom, fmtTo),
		"--pretty=format:'%H%s'",
	}

	// Convert excludes list into git grep commands
	if len(excludes) > 0 {
		fmtExcludes := make([]string, len(excludes))
		for i := range excludes {
			fmtExcludes[i] = fmt.Sprintf("--grep=%s", excludes[i])
		}
		fmtExcludes = append(fmtExcludes, "--invert-grep")

		// Append to original set of arguments
		args = append(args, fmtExcludes...)
	}

	log, err := Clean(Run(args...))
	if err != nil {
		return []LogEntry{}, err
	}

	if log == "" {
		return []LogEntry{}, nil
	}

	rows := strings.Split(log, "\n")
	les := make([]LogEntry, 0, len(rows))
	for _, r := range rows {
		les = append(les, LogEntry{
			Hash:       r[:40],
			AbbrevHash: r[:7],
			Message:    r[40:],
		})
	}

	return les, nil
}

// Staged retrieves a list of all files that are currently staged
func Staged() ([]string, error) {
	files, err := Clean(Run("diff", "--cached", "--name-only"))
	if err != nil {
		return []string{}, err
	}

	if files == "" {
		return []string{}, nil
	}

	return strings.Split(files, "\n"), nil
}

// Clean the output
func Clean(output string, err error) (string, error) {
	// Preserve multi-line output, but trim the trailing newline
	output = strings.TrimSuffix(strings.Replace(output, "'", "", -1), "\n")
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}
