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
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/gembaadvantage/codecommit-sign/pkg/translate"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// SCM is used for identifying the source code management tool used by the current
// git repository
type SCM string

const (
	GitHub       SCM = "GitHub"
	GitLab       SCM = "GitLab"
	CodeCommit   SCM = "CodeCommit"
	Gitea        SCM = "Gitea"
	Unrecognised SCM = "Unrecognised"
)

// CommitDetails contains mandatory details about a specific git commit
type CommitDetails struct {
	Message string
	Author  string
	Email   string
}

// AuthorDetails contains details about a git commit author
type AuthorDetails struct {
	Name  string
	Email string
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

// Repository contains details about a specific repository
type Repository struct {
	Origin string
	Owner  string
	Name   string
	Host   string
	Path   string
}

// String prints out a user friendly string representation
func (c CommitDetails) String() string {
	return fmt.Sprintf("%s <%s>\n%s", c.Author, c.Email, c.Message)
}

// Run executes a git command and returns its output or errors
func Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

// IsInstalled identifies whether git is installed under the current $PATH
func IsInstalled() bool {
	_, err := Run("--version")
	return err == nil
}

// IsRepo identifies whether the current working directory is a recognised
// git repository
func IsRepo() bool {
	out, err := Run("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// IsShallow identifies if the current repository was created through a shallow clone
func IsShallow() bool {
	out, err := Run("rev-parse", "--is-shallow-repository")
	return err == nil && strings.TrimSpace(out) == "true"
}

// CheckDirty identifies if the current repository is dirty through the presence of
// un-committed and/or un-staged changes and returns a list of those files
func CheckDirty() (string, error) {
	out, err := Clean(Run("status", "--porcelain"))
	if out != "" || err != nil {
		return out, err
	}
	return "", nil
}

// IsDeatched identifies if the current repository is detached from its HEAD
func IsDetached() bool {
	out, err := Run("branch", "--show-current")
	return err == nil && strings.TrimSpace(out) == ""
}

// Remote retrieves details about the remote origin of a repository
func Remote() (Repository, error) {
	remURL, err := Clean(Run("ls-remote", "--get-url"))
	if err != nil {
		return Repository{}, errors.New("no remote origin detected")
	}

	origin := remURL

	// Strip off any trailing .git suffix
	rem := strings.TrimSuffix(remURL, ".git")

	// Special use case for CodeCommit as that prefixes SSH URLs with ssh://
	rem = strings.TrimPrefix(rem, "ssh://")

	// Detect and translate a CodeCommit GRC URL into its HTTPS counterpart
	if strings.HasPrefix(rem, "codecommit:") {
		// Translate a codecommit GRC URL into its HTTPS counterpart
		if rem, err = translate.FromGRC(rem); err != nil {
			return Repository{}, err
		}

		origin = rem
	}

	if strings.HasPrefix(rem, "git@") {
		// Sanitise any SSH based URL to ensure it is parseable
		rem = strings.TrimPrefix(rem, "git@")
		rem = strings.Replace(rem, ":", "/", 1)
	}

	u, err := url.Parse(rem)
	if err != nil {
		return Repository{}, err
	}

	// Split into parts
	p := strings.Split(u.Path, "/")

	if len(p) < 3 {
		// This could be a custom Git server that doesn't follow the expected pattern.
		// Don't fail, but return the raw origin for custom parsing
		return Repository{Origin: origin}, nil
	}

	// If the repository has a HTTP(S) origin, the host will have been correctly identified
	host := u.Host
	if host == "" {
		// For other schemes, assume the host is contained within the first part
		host = p[0]
	}
	owner := p[1]
	name := p[len(p)-1]
	path := strings.Join(p[1:], "/")

	if strings.Contains(host, "codecommit") {
		// No concept of an owner with CodeCommit repositories
		owner = ""
		path = name
	}

	return Repository{
		Origin: origin,
		Owner:  owner,
		Name:   name,
		Path:   path,
		Host:   host,
	}, nil
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
	return retrieveTags([]string{"--sort=-creatordate"})
}

func retrieveTags(sort []string) []TagEntry {
	// Git can only perform basic pattern matching, so attempt a crude filtering for
	// semantic versions using the pattern *.*.* (major.minor.patch)
	args := []string{
		"for-each-ref",
		"refs/tags/*.*.*",
		`--format='%(creatordate:short),%(refname:short)'`,
	}
	args = append(args, sort...)

	tags, err := Clean(Run(args...))
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

		// Parse the ref to ensure it is a semantic version. If not, omit it from the identified tags
		v, err := semver.Parse(p[1])
		if err != nil {
			continue
		}

		ents = append(ents, TagEntry{
			Created: p[0],
			Ref:     v.Raw,
		})
	}

	return ents
}

// LatestTag retrieves the latest tag within the repository
func LatestTag() TagEntry {
	tags := retrieveTags([]string{"--sort=-creatordate", "--sort=-v:refname"})
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

// Log retrieves a log containing the commit history of a repository.
// If a tag is provided, the log will be generated from that tag to
// the current HEAD of the repository
func Log(tag string) (string, error) {
	if tag == "" {
		return commitLog("HEAD")
	}

	return commitLog(fmt.Sprintf("tags/%s..HEAD", tag))
}

func commitLog(srch string) (string, error) {
	out, err := Clean(Run("log", "--no-decorate", "--no-color", srch))
	if err != nil {
		return "", err
	}

	return out, nil
}

// Tag will create a lightweight tag against the repository
func Tag(tag string) error {
	if _, err := Clean(Run("tag", "-f", tag)); err != nil {
		return err
	}

	return nil
}

// AnnotatedTag will create an annotated tag against the repository
func AnnotatedTag(tag string, cd CommitDetails) error {
	args := []string{
		"-c",
		fmt.Sprintf("user.name='%s'", cd.Author),
		"-c",
		fmt.Sprintf("user.email='%s'", cd.Email),
		"tag",
		"-a",
		tag,
		"-f",
		"-m",
		cd.Message,
	}

	if _, err := Clean(Run(args...)); err != nil {
		return err
	}

	return nil
}

// PushTag attempts to push a newly created tag to the configured origin
func PushTag(tag string, options []string) error {
	args := []string{"push"}

	// Append each push option in the following format --push-option=<option>
	if len(options) > 0 {
		for _, opt := range options {
			args = append(args, fmt.Sprintf("--push-option=%s", opt))
		}
	}

	args = append(args, []string{"origin", tag}...)

	if _, err := Clean(Run(args...)); err != nil {
		return err
	}

	return nil
}

// Author attempts to retrieve details about the commit author directly
// from git config
func Author() AuthorDetails {
	name, _ := Clean(Run("config", "user.name"))
	email, _ := Clean(Run("config", "user.email"))

	return AuthorDetails{
		Name:  name,
		Email: email,
	}
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

	// If GPG commit signing is enabled, append the -S flag to the args
	if ConfigExists("commit.gpgsign", "true") {
		args = append(args, "-S")
	}

	if _, err := Clean(Run(args...)); err != nil {
		return err
	}

	return nil
}

// ConfigExists checks whether a given property is set within the local git
// config file of the repository
func ConfigExists(key, value string) bool {
	out, err := Clean(Run("config", "--get", key))
	if err != nil {
		return false
	}

	return out == value
}

// ConfigSet will attempt to set a series of git properties within the
// local config of the git repository
func ConfigSet(values map[string]string) error {
	for k, v := range values {
		if _, err := Run("config", "--add", k, v); err != nil {
			return err
		}
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
func Push(options []string) error {
	// Inspect the repo for an origin. If no origin exists, then skip the push
	if _, err := Clean(Run("remote", "show", "origin")); err != nil {
		return nil
	}

	branch, err := Clean(Run("rev-parse", "--abbrev-ref", "HEAD"))
	if err != nil {
		return err
	}

	args := []string{"push"}

	// Append each push option in the following format --push-option=<option>
	if len(options) > 0 {
		for _, opt := range options {
			args = append(args, fmt.Sprintf("--push-option=%s", opt))
		}
	}

	args = append(args, []string{"origin", branch}...)

	if _, err := Clean(Run(args...)); err != nil {
		return err
	}

	return nil
}

// LogBetween retrieves all log entries between two points of time within the
// git history of the repository. Supports tags and specific git hashes as its
// reference points. From must always be the closest point to HEAD
func LogBetween(from, to string) ([]LogEntry, error) {
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
		"--pretty=oneline",
		"--no-decorate",
		"--no-color",
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
		// Always a single whitespace between the hash and message
		hash, message, _ := strings.Cut(r, " ")
		les = append(les, LogEntry{
			Hash:       hash,
			AbbrevHash: hash[:7],
			Message:    message,
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
