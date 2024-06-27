package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	log := `fix: bug fix
docs: update docs
ci: update pipeline
feat: new feature`
	gittest.InitRepository(t, gittest.WithLog(log))

	tagCmd := newTagCmd(noChangesPushed(), os.Stdout)

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0", tags[0])
}

func TestTag_CurrentFlag(t *testing.T) {
	log := `(tag: v0.1.0) docs: updated docs
fix: bug fix`
	gittest.InitRepository(t, gittest.WithLog(log))

	var buf bytes.Buffer
	tagCmd := newTagCmd(noChangesPushed(), &buf)
	tagCmd.Cmd.SetArgs([]string{"--current"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "v0.1.0", buf.String())
}

func TestTag_NextFlag(t *testing.T) {
	log := `docs: updated docs
refactor!: breaking cli change
ci: update pipeline
fix: bug fix`
	gittest.InitRepository(t, gittest.WithLog(log))

	var buf bytes.Buffer
	tagCmd := newTagCmd(noChangesPushed(), &buf)
	tagCmd.Cmd.SetArgs([]string{"--next"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	assert.Len(t, tags, 0)
	assert.Equal(t, "v1.0.0", buf.String())
}

func TestTag_CurrentAndNextFlag(t *testing.T) {
	log := `fix: found another bug
(tag: v0.1.0) docs: updated docs
fix: bug fix`
	gittest.InitRepository(t, gittest.WithLog(log))

	var buf bytes.Buffer
	tagCmd := newTagCmd(noChangesPushed(), &buf)
	tagCmd.Cmd.SetArgs([]string{"--current", "--next"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0", tags[0])
	assert.Equal(t, "v0.1.0 v0.1.1", buf.String())
}

func TestTag_NoPrefix(t *testing.T) {
	log := `docs: update docs
fix: bug fix`
	gittest.InitRepository(t, gittest.WithLog(log))

	tagCmd := newTagCmd(noChangesPushed(), os.Stdout)
	tagCmd.Cmd.SetArgs([]string{"--no-prefix"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	assert.Len(t, tags, 1)
	assert.Equal(t, "0.0.1", tags[0])
}

func TestTag_PrereleaseFlag(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLog("feat: a new feature"))

	tagCmd := newTagCmd(noChangesPushed(), os.Stdout)
	tagCmd.Cmd.SetArgs([]string{"--prerelease", "-beta.1+12345"})

	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	tags := gittest.Tags(t)
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0-beta.1+12345", tags[0])
}

func TestTag_Hooks(t *testing.T) {
	gittest.InitRepository(t)
	configWithHooks(t)
	gittest.CommitEmpty(t, "feat: this is a new feature")

	tagCmd := newTagCmd(noChangesPushed(), os.Stdout)
	err := tagCmd.Cmd.Execute()
	require.NoError(t, err)

	require.Equal(t, 4, numHooksExecuted(t))
	assert.FileExists(t, BeforeFile)
	assert.FileExists(t, BeforeTagFile)
	assert.FileExists(t, AfterTagFile)
	assert.FileExists(t, AfterFile)
}
