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
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	assert.False(t, Task{}.Skip(&context.Context{}))
}

func TestString(t *testing.T) {
	assert.Equal(t, "checking git", Task{}.String())
}

func TestRun(t *testing.T) {
	git.InitRepo(t)

	err := Task{}.Run(&context.Context{})
	assert.NoError(t, err)
}

func TestRun_NoGit(t *testing.T) {
	// Just blast the PATH variable temporarily for this test
	t.Setenv("PATH", "")

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, "git is not currently installed under $PATH")
}

func TestRun_NotGitRepository(t *testing.T) {
	git.MkTmpDir(t)

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, "current working directory is not a git repository")
}

func TestRun_DetachedHead(t *testing.T) {
	git.InitRepo(t)
	h := git.EmptyCommit(t, "this is a test")

	// Checkout the returned hash to force a detached HEAD
	git.Run("checkout", h)

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, `uplift cannot reliably run when the repository is in a detached HEAD state. Some features
will not run as expected. To suppress this error, use the '--ignore-detached' flag, or 
set the required config.

For further details visit: https://upliftci.dev/faq/git-detached
`)
}

func TestRun_IgnoreDetachedHead(t *testing.T) {
	git.InitRepo(t)
	h := git.EmptyCommit(t, "this is a test")

	// Checkout the returned hash to force a detached HEAD
	git.Run("checkout", h)

	err := Task{}.Run(&context.Context{
		IgnoreDetached: true,
	})
	assert.NoError(t, err)
}

func TestRun_ShallowClone(t *testing.T) {
	git.InitShallowRepo(t)

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, `uplift cannot reliably run against a shallow clone of the repository. Some features may not 
work as expected. To suppress this error, use the '--ignore-shallow' flag, or set the 
required config.

For further details visit: https://upliftci.dev/faq/git-shallow
`)
}

func TestRun_IgnoreShallowClone(t *testing.T) {
	git.InitShallowRepo(t)

	err := Task{}.Run(&context.Context{
		IgnoreShallow: true,
	})
	assert.NoError(t, err)
}

func TestRun_Dirty(t *testing.T) {
	git.InitRepo(t)
	git.TouchFiles(t, "testing.go")

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, `uplift cannot reliably run if the repository is in a dirty state. Changes detected:
?? testing.go

Please check and resolve the status of these files before retrying. For further 
details visit: https://upliftci.dev/faq/git-dirty
`)
}
