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

package hook

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec_ShellCommands(t *testing.T) {
	git.MkTmpDir(t)

	cmds := []Command{
		{
			Operation: "echo -n 'JohnDoe' > out.txt",
		},
		{
			Operation: "sed --posix -i 's/Doe/Smith/g' out.txt",
		},
	}

	err := Exec(cmds)
	require.NoError(t, err)

	data, err := ioutil.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "JohnSmith", string(data))
}

func TestExec_ShellScripts(t *testing.T) {
	git.InitRepo(t)

	// Generate a shell script
	sh := `#!/bin/bash
LAST_COMMIT=$(git log -1 --pretty=format:'%B')
echo -n $LAST_COMMIT > out.txt`
	os.Mkdir("subfolder", 0755)
	ioutil.WriteFile("subfolder/last-commit.sh", []byte(sh), 0755)

	cmds := []Command{
		{
			Operation: "bash subfolder//last-commit.sh",
		},
	}

	err := Exec(cmds)
	require.NoError(t, err)

	data, err := ioutil.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "initialise repo", string(data))
}
