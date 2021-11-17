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

package changelog

import (
	"bytes"
	_ "embed"
	"errors"
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

	appendHeader = "## [Unreleased]\n\n"
)

var (
	//go:embed template/new.tpl
	newTpl string

	//go:embed template/append.tpl
	appendTpl string

	newTplBody    = template.Must(template.New("new").Parse(newTpl))
	appendTplBody = template.Must(template.New("append").Parse(appendTpl))

	// ErrNoAppendHeader is reported if a changelog is missing the expected append header
	ErrNoAppendHeader = errors.New("changelog missing supported append header")
)

type release struct {
	Tag     string
	Date    string
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
	if ctx.NextVersion.Raw == "" {
		log.Info("no release detected, skipping changelog")
		return nil
	}

	log.WithField("tag", ctx.NextVersion.Raw).Info("determine changes for release")
	ents, err := git.LogBetween(ctx.NextVersion.Raw, ctx.CurrentVersion.Raw)
	if err != nil {
		return err
	}

	// Package log entries into release ready for template generation
	rel := release{
		Tag:     ctx.NextVersion.Raw,
		Date:    time.Now().UTC().Format(ChangeDate),
		Changes: ents,
	}
	log.WithFields(log.Fields{
		"tag":     ctx.NextVersion.Raw,
		"date":    time.Now().UTC().Format(ChangeDate),
		"commits": len(rel.Changes),
	}).Info("changeset identified")

	if ctx.Debug {
		for _, c := range rel.Changes {
			log.WithFields(log.Fields{
				"hash":    c.AbbrevHash,
				"message": c.Message,
			}).Debug("commit")
		}
	}

	if ctx.DryRun {
		log.Info("skip writing to changelog in dry run mode")
		return nil
	}

	if noChangelogExists() {
		f, err := os.OpenFile(MarkdownFile, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		log.Debug("create new changelog in repository")
		return newTplBody.Execute(f, rel)
	}

	// When modifying an existing changelog, simply overwrite the existing file
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
	if err := appendTplBody.Execute(&buf, rel); err != nil {
		return err
	}

	apnd := strings.Replace(clStr, appendHeader, buf.String(), 1)

	log.Debug("append to existing changelog in repository")
	return ioutil.WriteFile(MarkdownFile, []byte(apnd), 0644)
}

func noChangelogExists() bool {
	_, err := os.Stat(MarkdownFile)
	return os.IsNotExist(err)
}
