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
	gittest.InitRepository(t)

	cmds := []string{
		"echo -n 'JohnDoe' > out.txt",
		"sed -i 's/Doe/Smith/g' out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "JohnSmith", string(data))
}

func TestExec_ShellScripts(t *testing.T) {
	log := "(tag: 1.0.0) feat: first release"
	gittest.InitRepository(t, gittest.WithLog(log))

	// Generate a shell script
	sh := `#!/bin/bash
	LATEST_TAG=$(git for-each-ref "refs/tags/*.*.*" --sort=-v:creatordate --format='%(refname:short)')
	echo -n $LATEST_TAG > out.txt`
	os.WriteFile("latest-tag.sh", []byte(sh), 0o755)

	err := Exec(context.Background(), []string{"./latest-tag.sh"}, ExecOptions{})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "1.0.0", string(data))
}
