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
	"io/ioutil"
	"os"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBump(t *testing.T) {
	git.InitRepo(t)
	testFileWithConfig(t, "test.txt", ".uplift.yml")
	git.EmptyCommitAndTag(t, "1.0.0", "feat: this is a new feature")

	bmpCmd := newBumpCmd(noChangesPushed(), os.Stdout)

	err := bmpCmd.Cmd.Execute()
	require.NoError(t, err)

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Equal(t, `version: 1.1.0
appVersion: 1.1.0`, string(actual))
}

func testFileWithConfig(t *testing.T, f string, cfg string) []byte {
	t.Helper()

	c := []byte(`version: 0.0.0
appVersion: 0.0.0`)
	err := ioutil.WriteFile(f, c, 0644)
	require.NoError(t, err)

	yml := fmt.Sprintf(`
bumps:
  - file: %s
    regex:
      - pattern: "version: $VERSION"
      - pattern: "appVersion: $VERSION"`, f)

	err = ioutil.WriteFile(cfg, []byte(yml), 0644)
	require.NoError(t, err)

	// Ensure files are committed to prevent dirty repository
	git.CommitFiles(t, f, cfg)
	return c
}

func TestBump_PrereleaseFlag(t *testing.T) {
	git.InitRepo(t)
	testFileWithConfig(t, "test.txt", ".uplift.yml")
	git.EmptyCommit(t, "feat: this is a new feature")

	bmpCmd := newBumpCmd(&globalOptions{}, os.Stdout)
	bmpCmd.Cmd.SetArgs([]string{"--prerelease", "-beta.1+12345"})

	err := bmpCmd.Cmd.Execute()
	require.NoError(t, err)

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Equal(t, `version: 0.1.0-beta.1+12345
appVersion: 0.1.0-beta.1+12345`, string(actual))
}
