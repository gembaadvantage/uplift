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

// VerboseLogger defines a logger that supports various different logging levels
// with text decoration. It is designed for logging out verbose information to aid
// the understanding of the application logic and shouldn't be used as the default
// logger for the application
type VerboseLogger struct {
	w io.Writer
}

// NewVerboseLogger creates a new verbose logger
func NewVerboseLogger(out io.Writer) ConsoleLogger {
	return VerboseLogger{
		w: out,
	}
}

// Out will log without any text decoration and line ending
func (l VerboseLogger) Out(s string, args ...interface{}) {
	fmt.Fprintf(l.w, "%s\n", fmt.Sprintf(s, args...))
}

// Success will log with success text decoration and line ending
func (l VerboseLogger) Success(s string, args ...interface{}) {
	fmt.Fprintf(l.w, "%s %s\n", GreenTick, fmt.Sprintf(s, args...))
}

// Info will log with information text decoration and line ending
func (l VerboseLogger) Info(s string, args ...interface{}) {
	// Double whitespace is needed to preserve formatting consistency
	fmt.Fprintf(l.w, "%s  %s\n", Information, fmt.Sprintf(s, args...))
}

// Warn will log with warning text decoration and line ending
func (l VerboseLogger) Warn(s string, args ...interface{}) {
	fmt.Fprintf(l.w, "%s %s\n", Warning, fmt.Sprintf(s, args...))
}
