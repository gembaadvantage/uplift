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
	"context"
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec_ShellCommands(t *testing.T) {
	// git.MkTmpDir(t)

	// cmds := []string{
	// 	"echo -n 'JohnDoe' > out.txt",
	// 	"sed --posix -i 's/Doe/Smith/g' out.txt",
	// }

	// err := Exec(context.Background(), cmds, ExecOptions{})
	// require.NoError(t, err)

	// data, err := os.ReadFile("out.txt")
	// require.NoError(t, err)

	// assert.Equal(t, "JohnSmith", string(data))
}

func TestExec_ShellScripts(t *testing.T) {
	gittest.InitRepository(t)

	// Generate a shell script
	sh := `#!/bin/bash
	LAST_COMMIT=$(git log -1 --pretty=format:'%B')
	echo -n $LAST_COMMIT > out.txt`
	os.Mkdir("subfolder", 0o755)
	os.WriteFile("subfolder/last-commit.sh", []byte(sh), 0o755)

	err := Exec(context.Background(), []string{"bash subfolder//last-commit.sh"}, ExecOptions{})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, gittest.InitialCommit, string(data))
}
