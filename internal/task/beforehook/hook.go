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

package beforehook

import (
	"os"
	"strings"

	"github.com/gembaadvantage/uplift/internal/context"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Task for executing any custom shell commands or scripts
// before uplift runs it release workflow
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "running before hooks"
}

// Skip running the task
func (t Task) Skip(ctx *context.Context) bool {
	return len(ctx.Config.Hooks.Before) == 0
}

// Run the task, executing any provided shell scripts or commands
func (t Task) Run(ctx *context.Context) error {
	for _, c := range ctx.Config.Hooks.Before {
		p, err := syntax.NewParser().Parse(strings.NewReader(c), "")
		if err != nil {
			return err
		}

		r, err := interp.New(interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
		if err != nil {
			return err
		}

		if err := r.Run(ctx.Context, p); err != nil {
			return err
		}
	}

	return nil
}
