package main

import (
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	bumpFile = `version: 0.0.0
appVersion: 0.0.0`

	bumpConfig = `bumps:
  - file: test.txt
    regex:
      - pattern: "version: $VERSION"
      - pattern: "appVersion: $VERSION"
`
)

func TestBump(t *testing.T) {
	log := `ci: update workflow"
docs: update docs
feat: a new feature
fix: a bug fix
(tag: 0.1.0) feat: this was the last feature`
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithCommittedFiles("test.txt", ".uplift.yml"),
		gittest.WithFileContent("test.txt", bumpFile, ".uplift.yml", bumpConfig))

	bmpCmd := newBumpCmd(noChangesPushed(), os.Stdout)

	err := bmpCmd.Cmd.Execute()
	require.NoError(t, err)

	actual, err := os.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Equal(t, `version: 0.2.0
appVersion: 0.2.0`, string(actual))
}

func TestBump_PrereleaseFlag(t *testing.T) {
	log := `docs: update docs
fix: fix bug
feat!: breaking change
feat: this is a new feature`
	gittest.InitRepository(t,
		gittest.WithLog(log),
		gittest.WithCommittedFiles("test.txt", ".uplift.yml"),
		gittest.WithFileContent("test.txt", bumpFile, ".uplift.yml", bumpConfig))

	bmpCmd := newBumpCmd(&globalOptions{}, os.Stdout)
	bmpCmd.Cmd.SetArgs([]string{"--prerelease", "-beta.1+12345"})

	err := bmpCmd.Cmd.Execute()
	require.NoError(t, err)

	actual, err := os.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Equal(t, `version: v1.0.0-beta.1+12345
appVersion: v1.0.0-beta.1+12345`, string(actual))
}

func TestBump_Hooks(t *testing.T) {
	gittest.InitRepository(t, gittest.WithLog("feat: this is a new feature"))
	configWithHooks(t)

	bmpCmd := newBumpCmd(&globalOptions{}, os.Stdout)
	err := bmpCmd.Cmd.Execute()
	require.NoError(t, err)

	require.Equal(t, 4, numHooksExecuted(t))
	assert.FileExists(t, BeforeFile)
	assert.FileExists(t, BeforeBumpFile)
	assert.FileExists(t, AfterBumpFile)
	assert.FileExists(t, AfterFile)
}
