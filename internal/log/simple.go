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

package log

import (
	"fmt"
	"io"
)

// SimpleLogger defines a logger that logs without any text decoration
// and only supports logging using standard output. Logging with any of the
// provided logging levels will simply be ignored
type SimpleLogger struct {
	w io.Writer
}

// NewSimpleLogger creates a new simple logger
func NewSimpleLogger(out io.Writer) ConsoleLogger {
	return SimpleLogger{
		w: out,
	}
}

// Out will log without any text decoration
func (l SimpleLogger) Out(s string, args ...interface{}) {
	fmt.Fprintf(l.w, s, args...)
}

// Success will not log anything
func (l SimpleLogger) Success(s string, args ...interface{}) {
}

// Info will not log anything
func (l SimpleLogger) Info(s string, args ...interface{}) {
}

// Warn will  not log anything
func (l SimpleLogger) Warn(s string, args ...interface{}) {
}
