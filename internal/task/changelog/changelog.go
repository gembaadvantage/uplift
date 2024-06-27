package changelog

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	git "github.com/purpleclay/gitz"
)

const (
	//  MarkdownFile defines the name of the changelog markdown file within a repository
	MarkdownFile = "CHANGELOG.md"

	// ChangeDate defines a date formatting used when logging a new change
	ChangeDate = "2006-01-02"

	appendHeader = "## Unreleased\n\n"
)

var (
	//go:embed template/new.tmpl
	newTpl string

	//go:embed template/append.tmpl
	appendTpl string

	//go:embed template/diff.tmpl
	diffTpl string

	funcs = template.FuncMap{
		"tpl": execTemplate,
	}

	newTplBody    = template.Must(template.New("new").Funcs(funcs).Parse(newTpl))
	appendTplBody = template.Must(template.New("append").Funcs(funcs).Parse(appendTpl))
	diffTplBody   = template.Must(template.New("diff").Funcs(funcs).Parse(diffTpl))

	// ErrNoAppendHeader is reported if a changelog is missing the expected append header
	ErrNoAppendHeader = errors.New("changelog missing supported append header")
)

type release struct {
	SCM     context.SCM
	Tag     tagEntry
	Changes []git.LogEntry
}

type tagEntry struct {
	Ref     string
	Created string
}

// Task that generates a changelog for the current repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "generating changelog"
}

// Skip running the task if no changelog is needed
func (t Task) Skip(ctx *context.Context) bool {
	if ctx.NoVersionChanged || ctx.SkipChangelog {
		return true
	}

	if ctx.Changelog.SkipPrerelease && ctx.NextVersion.Prerelease != "" {
		return true
	}

	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	if !ctx.Changelog.All && ctx.NextVersion.Raw == "" {
		log.Info("no release detected, skipping changelog")
		return nil
	}

	// To ensure the changelog is correctly generated during a release, the latest commit will
	// need to be temporarily tagged. After changelog generation, it will be removed to
	// ensure the release workflow behaves as expected. This will be a transparent operation
	// that cannot be invoked by the caller
	if ctx.Changelog.PreTag {
		log.WithField("tag", ctx.NextVersion.Raw).Info("pre-tagging latest commit for changelog creation")
		if _, err := ctx.GitClient.Tag(ctx.NextVersion.Raw, git.WithLocalOnly()); err != nil {
			return err
		}
		defer func() {
			log.Info("removing pre-tag after changelog creation")
			if _, err := ctx.GitClient.DeleteTag(ctx.NextVersion.Raw, git.WithLocalDelete()); err != nil {
				log.WithError(err).Error("failed to delete pre-tag")
			}
		}()
	}

	// Retrieve log entries based on the changelog expectations
	var rels []release
	var relErr error
	if ctx.Changelog.All {
		rels, relErr = changelogReleases(ctx)
	} else {
		rels, relErr = changelogRelease(ctx)
	}
	if relErr != nil {
		return relErr
	}

	// If there are no releases to changelog, simply return
	if len(rels) == 0 {
		return nil
	}

	if ctx.DryRun {
		log.Info("skip writing to changelog in dry run mode")
		return nil
	}

	if ctx.Changelog.Multiline {
		log.Info("formatting multiline messages for changelog")
		for i := range rels {
			for j := range rels[i].Changes {
				msg := rels[i].Changes[j].Message
				if ctx.Changelog.TrimHeader {
					startIdx := semver.FindStartIdx(msg)
					msg = msg[startIdx:]
				}
				msg = strings.ReplaceAll(msg, "\n", "\n  ")
				msg = strings.ReplaceAll(msg, "\n  \n", "\n\n")

				rels[i].Changes[j].Message = msg
			}
		}
	} else {
		log.Info("trim all commit messages to a single line")
		for i := range rels {
			for j := range rels[i].Changes {
				msg := rels[i].Changes[j].Message
				if ctx.Changelog.TrimHeader {
					startIdx := semver.FindStartIdx(msg)
					msg = msg[startIdx:]
					rels[i].Changes[j].Message = msg
				}
				if idx := strings.Index(msg, "\n"); idx > -1 {
					rels[i].Changes[j].Message = strings.TrimSpace(msg)
				}
			}
		}
	}

	if ctx.Changelog.DiffOnly {
		diff, err := diffChangelog(rels)
		if err != nil {
			return err
		}
		fmt.Fprint(ctx.Out, diff)
		return nil
	}

	var chgErr error
	if noChangelogExists() || ctx.Changelog.All {
		chgErr = newChangelog(rels)
	} else {
		chgErr = appendChangelog(rels)
	}
	if chgErr != nil {
		return chgErr
	}

	if ctx.NoStage {
		log.Info("skip staging of CHANGELOG.md")
		return nil
	}

	log.Debug("staging CHANGELOG.md")
	_, err := ctx.GitClient.Stage(git.WithPathSpecs(MarkdownFile))
	return err
}

func changelogRelease(ctx *context.Context) ([]release, error) {
	next := ctx.NextVersion.Raw
	prev := ctx.CurrentVersion.Raw

	log.WithField("tag", next).Info("determine changes for release")
	if ctx.Changelog.SkipPrerelease {
		// Retrieve all tags and filter out any that are prerelease versions
		tags, _ := ctx.GitClient.Tags(git.WithShellGlob("*.*.*"),
			git.WithSortBy(git.CreatorDateDesc, git.VersionDesc),
			git.WithFilters(func(tag string) bool {
				ver, err := semver.Parse(tag)
				if err != nil {
					return false
				}

				return ver.Prerelease == ""
			}),
			git.WithCount(2))

		if len(tags) == 1 {
			prev = ""
		} else {
			prev = tags[1]
		}
	}

	glog, err := ctx.GitClient.Log(git.WithRefRange(next, prev))
	if err != nil {
		return []release{}, err
	}

	ents := glog.Commits
	if len(ctx.Changelog.Include) > 0 {
		log.Info("cherry-picking commits based on include list")
		ents, err = includeCommits(ents, ctx.Changelog.Include)
		if err != nil {
			return []release{}, err
		}
	}

	if len(ctx.Changelog.Exclude) > 0 {
		log.Info("removing commits based on exclude list")
		ents, err = excludeCommits(ents, ctx.Changelog.Exclude)
		if err != nil {
			return []release{}, err
		}
	}

	if len(ents) == 0 {
		log.WithFields(log.Fields{
			"tag":  next,
			"prev": prev,
		}).Info("no log entries between tags")
		return []release{}, nil
	}

	log.WithFields(log.Fields{
		"tag":     next,
		"date":    time.Now().UTC().Format(ChangeDate),
		"commits": len(ents),
	}).Info("changeset identified")

	if ctx.Debug {
		for _, c := range ents {
			log.WithFields(log.Fields{
				"hash":    c.AbbrevHash,
				"message": c.Message,
			}).Debug("commit")
		}
	}

	if ctx.Changelog.Sort == "asc" {
		reverse(ents)
	}

	tagDetails, _ := ctx.GitClient.ShowTags(ctx.NextVersion.Raw)
	return []release{
		{
			SCM:     ctx.SCM,
			Tag:     extractTagEntry(tagDetails[ctx.NextVersion.Raw]),
			Changes: ents,
		},
	}, nil
}

func extractTagEntry(dets git.TagDetails) tagEntry {
	created := dets.Commit.CommitterDate.Format(ChangeDate)
	if dets.Annotation != nil {
		created = dets.Annotation.TaggerDate.Format(ChangeDate)
	}

	return tagEntry{
		Ref:     dets.Ref,
		Created: created,
	}
}

func changelogReleases(ctx *context.Context) ([]release, error) {
	tags, err := ctx.GitClient.Tags(git.WithShellGlob("*.*.*"),
		git.WithSortBy(git.CreatorDateDesc, git.VersionDesc),
		git.WithFilters(func(tag string) bool {
			if !ctx.Changelog.SkipPrerelease {
				return true
			}

			ver, err := semver.Parse(tag)
			if err != nil {
				return false
			}

			return ver.Prerelease == ""
		}))
	if err != nil {
		return []release{}, nil
	}

	if len(tags) == 0 {
		log.Info("no tags found within repository")
		return []release{}, nil
	}

	rels := make([]release, 0, len(tags))
	for i := 0; i < len(tags); i++ {
		nextTag := ""
		if i+1 < len(tags) {
			nextTag = tags[i+1]
		}

		tagDetails, _ := ctx.GitClient.ShowTags(tags[i])
		tag := extractTagEntry(tagDetails[tags[i]])

		log.WithField("tag", tags[i]).Info("determine changes for release")
		glog, err := ctx.GitClient.Log(git.WithRefRange(tag.Ref, nextTag))
		if err != nil {
			return []release{}, err
		}

		ents := glog.Commits
		if len(ctx.Changelog.Include) > 0 {
			log.Info("cherry-picking commits based on include list")
			ents, err = includeCommits(ents, ctx.Changelog.Include)
			if err != nil {
				return []release{}, err
			}
		}

		if len(ctx.Changelog.Exclude) > 0 {
			log.Info("removing commits based on exclude list")
			ents, err = excludeCommits(ents, ctx.Changelog.Exclude)
			if err != nil {
				return []release{}, err
			}
		}

		if len(ents) == 0 {
			log.WithFields(log.Fields{
				"tag":  tag.Ref,
				"prev": nextTag,
			}).Info("no log entries between tags")
		} else {
			log.WithFields(log.Fields{
				"tag":     tag.Ref,
				"date":    tag.Created,
				"commits": len(ents),
			}).Info("changeset identified")

			if ctx.Debug {
				for _, c := range ents {
					log.WithFields(log.Fields{
						"hash":    c.AbbrevHash,
						"message": c.Message,
					}).Debug("commit")
				}
			}
		}

		if ctx.Changelog.Sort == "asc" {
			reverse(ents)
		}

		rels = append(rels, release{
			SCM:     ctx.SCM,
			Tag:     tag,
			Changes: ents,
		})
	}

	return rels, nil
}

func noChangelogExists() bool {
	_, err := os.Stat(MarkdownFile)
	return os.IsNotExist(err)
}

func diffChangelog(rels []release) (string, error) {
	var buf bytes.Buffer
	if err := diffTplBody.Execute(&buf, rels); err != nil {
		return "", err
	}

	// Trim leading whitespace
	diff := buf.String()
	return diff[1:], nil
}

func newChangelog(rels []release) error {
	f, err := os.OpenFile(MarkdownFile, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debug("create new changelog in repository")
	return newTplBody.Execute(f, rels)
}

func appendChangelog(rels []release) error {
	cl, err := os.ReadFile(MarkdownFile)
	if err != nil {
		return err
	}

	// Appending is only possible if token token exists
	clStr := string(cl)
	if idx := strings.Index(clStr, appendHeader); idx == -1 {
		return ErrNoAppendHeader
	}

	var buf bytes.Buffer
	if err := appendTplBody.Execute(&buf, rels); err != nil {
		return err
	}

	apnd := strings.Replace(clStr, appendHeader, buf.String(), 1)

	log.Debug("append to existing changelog in repository")
	return os.WriteFile(MarkdownFile, []byte(apnd), 0o644)
}

func execTemplate(tmpl string, v interface{}) string {
	t, _ := template.New("dynamic").Parse(tmpl)

	var buf bytes.Buffer
	t.Execute(&buf, v)

	return buf.String()
}

func reverse(ents []git.LogEntry) {
	for i, j := 0, len(ents)-1; i < j; {
		ents[i], ents[j] = ents[j], ents[i]
		i++
		j--
	}
}

func includeCommits(commits []git.LogEntry, regexes []string) ([]git.LogEntry, error) {
	filtered := []git.LogEntry{}
	for _, regex := range regexes {
		includeRgx, err := regexp.Compile(regex)
		if err != nil {
			return filtered, err
		}

		// Iterate over the entire list of log entries for each regex and
		// append any match to the filtered list
		for _, commit := range commits {
			if includeRgx.MatchString(commit.Message) {
				filtered = append(filtered, commit)
			}
		}
	}

	return filtered, nil
}

func excludeCommits(commits []git.LogEntry, regexes []string) ([]git.LogEntry, error) {
	filtered := commits
	for _, regex := range regexes {
		excludeRgx, err := regexp.Compile(regex)
		if err != nil {
			return filtered, err
		}

		// Repeat over the filtered list for every exclude, compressing the list
		// of log entries on each iteration
		filterPass := []git.LogEntry{}
		for _, commit := range filtered {
			if !excludeRgx.MatchString(commit.Message) {
				filterPass = append(filterPass, commit)
			}
		}
		filtered = filterPass
	}

	return filtered, nil
}
