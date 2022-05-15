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

package hook

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/apex/log"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Command ...
type Command struct {
	Operation string
	Debug     bool
	DryRun    bool
}

// Exec ...
func Exec(cmds []Command) error {
	for _, c := range cmds {
		log.WithField("hook", c).Info("running")
		if c.DryRun {
			continue
		}

		p, err := syntax.NewParser().Parse(strings.NewReader(c.Operation), "")
		if err != nil {
			return err
		}

		// Discard all output from commands and scripts unless in debug mode
		out := io.Discard
		if c.Debug {
			// Stderr is used by apex for logging, stdout is reserved for capturing output
			out = os.Stderr
		}

		r, err := interp.New(
			interp.StdIO(os.Stdin, out, os.Stderr),
			interp.OpenHandler(openHandler),
		)
		if err != nil {
			return err
		}

		if err := r.Run(context.Background(), p); err != nil {
			return err
		}
	}

	return nil
}

func openHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return DevNull{}, nil
	}

	return interp.DefaultOpenHandler()(ctx, path, flag, perm)
}
