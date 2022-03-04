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

package changelog

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
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
	Tag     git.TagEntry
	Changes []git.LogEntry
}

// Task that generates a changelog for the current repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "changelog"
}

// Skip running the task if no changelog is needed
func (t Task) Skip(ctx *context.Context) bool {
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
		if err := git.Tag(ctx.NextVersion.Raw); err != nil {
			return err
		}
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

	log.Debug("staging CHANGELOG.md")
	return git.Stage(MarkdownFile)
}

func changelogRelease(ctx *context.Context) ([]release, error) {
	log.WithField("tag", ctx.NextVersion.Raw).Info("determine changes for release")
	ents, err := git.LogBetween(ctx.NextVersion.Raw, ctx.CurrentVersion.Raw, ctx.Changelog.Exclude)
	if err != nil {
		return []release{}, err
	}

	if len(ents) == 0 {
		log.WithFields(log.Fields{
			"tag":  ctx.NextVersion.Raw,
			"prev": ctx.CurrentVersion.Raw,
		}).Info("no log entries between tags")
		return []release{}, nil
	}

	log.WithFields(log.Fields{
		"tag":     ctx.NextVersion.Raw,
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

	return []release{
		{
			SCM:     ctx.SCM,
			Tag:     git.DescribeTag(ctx.NextVersion.Raw),
			Changes: ents,
		},
	}, nil
}

func changelogReleases(ctx *context.Context) ([]release, error) {
	tags := git.AllTags()
	if len(tags) == 0 {
		log.Info("no tags found within repository")
		return []release{}, nil
	}

	rels := make([]release, 0, len(tags))
	for i := 0; i < len(tags); i++ {
		nextTag := ""
		if i+1 < len(tags) {
			nextTag = tags[i+1].Ref
		}

		log.WithField("tag", tags[i].Ref).Info("determine changes for release")
		ents, err := git.LogBetween(tags[i].Ref, nextTag, ctx.Changelog.Exclude)
		if err != nil {
			return []release{}, err
		}

		if len(ents) == 0 {
			log.WithFields(log.Fields{
				"tag":  tags[i].Ref,
				"prev": nextTag,
			}).Info("no log entries between tags")
		} else {
			log.WithFields(log.Fields{
				"tag":     tags[i].Ref,
				"date":    tags[i].Created,
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
			Tag:     tags[i],
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
	f, err := os.OpenFile(MarkdownFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debug("create new changelog in repository")
	return newTplBody.Execute(f, rels)
}

func appendChangelog(rels []release) error {
	cl, err := ioutil.ReadFile(MarkdownFile)
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
	return ioutil.WriteFile(MarkdownFile, []byte(apnd), 0644)
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
