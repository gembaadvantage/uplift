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
	"io/ioutil"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec_DryRun(t *testing.T) {
	git.MkTmpDir(t)

	err := Exec(context.Background(), []string{"touch out.txt"}, ExecOptions{DryRun: true})
	require.NoError(t, err)

	assert.NoFileExists(t, "out.txt")
}

func TestExec_InjectEnvVars(t *testing.T) {
	git.MkTmpDir(t)

	env := []string{"ONE=1", "TWO=2"}

	sh := `#!/bin/bash
echo -n "ONE=$ONE TWO=$TWO" > out.txt`
	ioutil.WriteFile("print-env.sh", []byte(sh), 0o755)

	cmds := []string{
		"./print-env.sh",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: env})
	require.NoError(t, err)

	data, err := ioutil.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "ONE=1 TWO=2", string(data))
}

func TestExec_MergesEnvVars(t *testing.T) {
	git.MkTmpDir(t)

	env := []string{"ONE=1"}

	sh := `#!/bin/bash
printenv > out.txt`
	ioutil.WriteFile("list-env.sh", []byte(sh), 0o755)

	cmds := []string{
		"./list-env.sh",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: env})
	require.NoError(t, err)

	data, err := ioutil.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Contains(t, string(data), "PATH=")
	assert.Contains(t, string(data), "PWD=")
	assert.Contains(t, string(data), "HOME=")
	assert.Contains(t, string(data), "ONE=")
}
