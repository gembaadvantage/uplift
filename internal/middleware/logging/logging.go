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

package logging

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/middleware"
)

const (
	// DefaultPadding ensures all titles are indented by a set number of spaces
	DefaultPadding = 3

	// PrettyPadding ensures all other logging is indented twice the size of the default padding
	PrettyPadding = DefaultPadding * 2
)

// Log executes the given action and ensures the output is pretty printed.
// The title will always be followed by any indented output from the
// action itself
func Log(title string, act middleware.Action) middleware.Action {
	return func(ctx *context.Context) error {
		defer func() {
			cli.Default.Padding = DefaultPadding
		}()

		cli.Default.Padding = DefaultPadding
		log.Info(title)
		cli.Default.Padding = PrettyPadding

		return act(ctx)
	}
}
