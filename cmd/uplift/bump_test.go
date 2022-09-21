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
	"fmt"
	"os"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBump(t *testing.T) {
	taggedRepo(t, "0.1.0", "feat: this was the last feature")
	testFileWithConfig(t, "test.txt", ".uplift.yml")
	git.EmptyCommits(t, "ci: update workflow", "docs: update docs", "feat: a new feature", "fix: a bug fix")

	bmpCmd := newBumpCmd(noChangesPushed(), os.Stdout)

	err := bmpCmd.Cmd.Execute()
	require.NoError(t, err)

	actual, err := os.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Equal(t, `version: 0.2.0
appVersion: 0.2.0`, string(actual))
}

func testFileWithConfig(t *testing.T, f string, cfg string) []byte {
	t.Helper()

	c := []byte(`version: 0.0.0
appVersion: 0.0.0`)
	err := os.WriteFile(f, c, 0o644)
	require.NoError(t, err)

	yml := fmt.Sprintf(`
bumps:
  - file: %s
    regex:
      - pattern: "version: $VERSION"
      - pattern: "appVersion: $VERSION"`, f)

	err = os.WriteFile(cfg, []byte(yml), 0o644)
	require.NoError(t, err)

	// Ensure files are committed to prevent dirty repository
	git.CommitFiles(t, f, cfg)
	return c
}

func TestBump_PrereleaseFlag(t *testing.T) {
	untaggedRepo(t, "docs: update docs", "fix: fix bug", "feat!: breaking change")
	testFileWithConfig(t, "test.txt", ".uplift.yml")
	git.EmptyCommit(t, "feat: this is a new feature")

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
	untaggedRepo(t, "feat: this is a new feature")
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
