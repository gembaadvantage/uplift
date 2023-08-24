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

package gitcheck

import (
	"errors"
	"fmt"
)

// ErrDirty is raised when a git repository has un-committed and/or un-staged changes
type ErrDirty struct {
	status string
}

// Error returns a formatted message of the current error
func (e ErrDirty) Error() string {
	return fmt.Sprintf(`uplift cannot reliably run if the repository is in a dirty state. Changes detected:
%s

Please check and resolve the status of these files before retrying. For further
details visit: https://upliftci.dev/faq/gitdirty
`, e.status)
}

// ErrDetachedHead is raised if the git repository is in a detached state, which prevents
// some of the functionality of uplift from working correctly
type ErrDetachedHead struct{}

// Error returns a formatted message of the current error
func (e ErrDetachedHead) Error() string {
	return `uplift cannot reliably run when the repository is in a detached HEAD state. Some features
will not run as expected. To suppress this error, use the '--ignore-detached' flag, or
set the required config.

For further details visit: https://upliftci.dev/faq/gitdetached
`
}

// ErrShallowClone is raised if the git repository is from a shallow clone, which prevents
// some of the functionality of uplift from working correctly
type ErrShallowClone struct{}

// Error returns a formatted message of the current error
func (e ErrShallowClone) Error() string {
	return `uplift cannot reliably run against a shallow clone of the repository. Some features may not
work as expected. To suppress this error, use the '--ignore-shallow' flag, or set the
required config.

For further details visit: https://upliftci.dev/faq/gitshallow
`
}

var (
	// ErrGitMissing is raised if git is not detected on the current $PATH
	ErrGitMissing = errors.New("git is not currently installed under $PATH")

	// ErrNoRepository is raised if the current directory is not a recognised git repository
	ErrNoRepository = errors.New("current working directory is not a git repository")
)
