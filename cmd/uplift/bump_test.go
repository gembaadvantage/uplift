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
	gittest.InitRepository(t, gittest.WithLog(log))

	// TODO:
	// gittest.StagedFile(t, "", "") > TempFile and StageFile under the covers
	// gittest.WithCommittedFiles("", "", "", "") > file is automatically committed at the end

	gittest.TempFile(t, "test.txt", bumpFile)
	gittest.StageFile(t, "test.txt")
	gittest.TempFile(t, ".uplift.yml", bumpConfig)
	gittest.StageFile(t, ".uplift.yml")
	gittest.Commit(t, "chore: added files")

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
	gittest.InitRepository(t, gittest.WithLog(log))
	gittest.TempFile(t, "test.txt", bumpFile)
	gittest.StageFile(t, "test.txt")
	gittest.TempFile(t, ".uplift.yml", bumpConfig)
	gittest.StageFile(t, ".uplift.yml")
	gittest.Commit(t, "chore: added files")

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
