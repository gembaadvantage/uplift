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
	_ "embed"
	"os"
	"text/template"
	"time"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
)

//go:embed template/new.tpl
var new string

// Release ...
type Release struct {
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
	return ctx.NoChangelog
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	if ctx.NextVersion.Raw == "" {
		// TODO: log
		return nil
	}

	ents, err := git.LogBetween(ctx.NextVersion.Raw, ctx.CurrentVersion.Raw)
	if err != nil {
		return err
	}

	f, err := os.OpenFile("CHANGELOG.md", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl := template.Must(template.New("changelog").Parse(new))
	return tmpl.Execute(f, Release{
		Tag:     ctx.NextVersion.Raw,
		Date:    time.Now().UTC().Format("2006-01-02"),
		Changes: ents,
	})
}
