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
	"os"
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkip(t *testing.T) {
	assert.False(t, Task{}.Skip(&context.Context{}))
}

func TestString(t *testing.T) {
	assert.Equal(t, "checking git", Task{}.String())
}

func TestRun(t *testing.T) {
	gittest.InitRepository(t)

	err := Task{}.Run(&context.Context{})
	assert.NoError(t, err)
}

// TODO: this should be carried out early on in cobra

func TestRun_NotGitRepository(t *testing.T) {
	mkTmpDir(t)

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, "current working directory is not a git repository")
}

func mkTmpDir(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	current, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))

	t.Cleanup(func() {
		require.NoError(t, os.Chdir(current))
	})
}

func TestRun_DetachedHead(t *testing.T) {
	log := "(main) testing a detached head"
	gittest.InitRepository(t, gittest.WithLog(log))
	commit := gittest.LastCommit(t)
	gittest.Checkout(t, commit.Hash)

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, `uplift cannot reliably run when the repository is in a detached HEAD state. Some features
will not run as expected. To suppress this error, use the '--ignore-detached' flag, or
set the required config.

For further details visit: https://upliftci.dev/faq/gitdetached
`)
}

func TestRun_IgnoreDetachedHead(t *testing.T) {
	log := "(main) testing a detached head"
	gittest.InitRepository(t, gittest.WithLog(log))
	commit := gittest.LastCommit(t)
	gittest.Checkout(t, commit.Hash)

	err := Task{}.Run(&context.Context{
		IgnoreDetached: true,
	})
	assert.NoError(t, err)
}

func TestRun_ShallowClone(t *testing.T) {
	log := "(main) testing a shallow clone"
	gittest.InitRepository(t, gittest.WithLog(log), gittest.WithCloneDepth(1))

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, `uplift cannot reliably run against a shallow clone of the repository. Some features may not
work as expected. To suppress this error, use the '--ignore-shallow' flag, or set the
required config.

For further details visit: https://upliftci.dev/faq/gitshallow
`)
}

func TestRun_IgnoreShallowClone(t *testing.T) {
	log := "(main) testing a shallow clone"
	gittest.InitRepository(t, gittest.WithLog(log), gittest.WithCloneDepth(1))

	err := Task{}.Run(&context.Context{
		IgnoreShallow: true,
	})
	assert.NoError(t, err)
}

func TestRun_Dirty(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("testing.go"))

	err := Task{}.Run(&context.Context{})
	assert.EqualError(t, err, `uplift cannot reliably run if the repository is in a dirty state. Changes detected:
?? testing.go

Please check and resolve the status of these files before retrying. For further
details visit: https://upliftci.dev/faq/gitdirty
`)
}

func TestRun_DirtyWithConfiguredFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("testing.go"))

	err := Task{}.Run(&context.Context{
		Config: config.Uplift{
			Git: &config.Git{
				DirtyFiles: []string{"testing.go"},
			},
		},
	})
	assert.NoError(t, err)
}
