package gitcheck

import (
	"os"
	"testing"

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
[?? testing.go]

Please check and resolve the status of these files before retrying. For further
details visit: https://upliftci.dev/faq/gitdirty
`)
}

func TestRun_IncludeArtifactsWithConfiguredFiles(t *testing.T) {
	gittest.InitRepository(t, gittest.WithFiles("testing.go"))

	err := Task{}.Run(&context.Context{
		IncludeArtifacts: []string{"testing.go"},
	})
	assert.NoError(t, err)
}
