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

	env := []string{"VARIABLE=VALUE", "ANOTHER_VARIABLE=ANOTHER VALUE"}

	cmds := []string{
		"echo -n VALUE1=$VARIABLE VALUE2=$ANOTHER_VARIABLE > out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: env})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "VALUE1=VALUE VALUE2=ANOTHER VALUE", string(data))
}

func TestExec_MergesEnvVars(t *testing.T) {
	git.MkTmpDir(t)

	env := []string{"TESTING=123"}

	cmds := []string{
		"printenv > out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: env})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Contains(t, string(data), "PATH=")
	assert.Contains(t, string(data), "PWD=")
	assert.Contains(t, string(data), "HOME=")
	assert.Contains(t, string(data), "TESTING=123")
}

func TestExec_VarsWithWhitespace(t *testing.T) {
	git.MkTmpDir(t)

	env := []string{"ONE = 1", "TWO= 2", "THREE    =    3"}

	cmds := []string{
		"echo -n $ONE $TWO $THREE > out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: env})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "1 2 3", string(data))
}

func TestExec_LoadDotEnvFiles(t *testing.T) {
	git.MkTmpDir(t)

	dotenv1 := `ONE=1
TWO   =   2`
	os.WriteFile(".env", []byte(dotenv1), 0o600)

	os.Mkdir("custom", 0o755)
	dotenv2 := "THREE=    3"
	os.WriteFile("custom/another.env", []byte(dotenv2), 0o600)

	cmds := []string{
		"echo -n $ONE $TWO $THREE > out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{Env: []string{".env", "custom/another.env"}})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "1 2 3", string(data))
}

func TestExec_FailsOnInvalidDotEnvFile(t *testing.T) {
	git.MkTmpDir(t)

	dotenv := "INVALID"
	os.WriteFile(".env", []byte(dotenv), 0o600)

	err := Exec(context.Background(), []string{}, ExecOptions{Env: []string{".env"}})
	require.Error(t, err)
}

func TestExec_FailsIfDotEnvFileNotFound(t *testing.T) {
	git.MkTmpDir(t)

	err := Exec(context.Background(), []string{}, ExecOptions{Env: []string{"does-not-exist.env"}})
	require.Error(t, err)
}
