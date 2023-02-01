/*
Copyright (c) 2023 Gemba Advantage

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

package task

import (
	"fmt"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gembaadvantage/uplift/internal/context"
)

const (
	// DefaultPadding ensures all titles are indented by a set number of spaces
	DefaultPadding = 3

	// PrettyPadding ensures all other logging is indented twice the size of the default padding
	PrettyPadding = DefaultPadding * 2
)

// Runner defines a way of running a task. A task can either be run as a
// standalone operation or chained into a series of consecutive operations
type Runner interface {
	fmt.Stringer

	// Run the task. A context is provided allowing state between tasks to
	// be shared. Useful if multiple tasks are executed in order
	Run(ctx *context.Context) error

	// Skip running of the task based on the current context state
	Skip(ctx *context.Context) bool
}

// Execute a series of tasks, providing the [context.Context] to each. Before executing
// a task, a precondition check is performed, identifying if the task should be skipped
// or not. Tasks that are skipped, will automatically have [skipped] appended to their
// task name. Execution will be aborted upon the first encountered error
func Execute(ctx *context.Context, tasks []Runner) error {
	for _, t := range tasks {
		defer func() {
			// Ensure padding is automatically reset
			cli.Default.Padding = DefaultPadding
		}()

		cli.Default.Padding = DefaultPadding

		if t.Skip(ctx) {
			log.Debug(fmt.Sprintf("(skipped) %s", t.String()))
		} else {
			log.Info(t.String())
			cli.Default.Padding = PrettyPadding
			if err := t.Run(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}
