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

package gitcommit

import (
	"io/ioutil"
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "committing changes", Task{}.String())
}

func TestSkip(t *testing.T) {
	tests := []struct {
		name string
		ctx  *context.Context
	}{
		{
			name: "DryRun",
			ctx: &context.Context{
				DryRun: true,
			},
		},
		{
			name: "NoVersionChanged",
			ctx: &context.Context{
				NoVersionChanged: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, Task{}.Skip(tt.ctx))
		})
	}
}

func TestRun(t *testing.T) {
	git.InitRepo(t)
	trackFile(t, "test.txt")

	err := Task{}.Run(&context.Context{
		CommitDetails: git.CommitDetails{
			Author:  "uplift",
			Email:   "uplift@test.com",
			Message: "test commit",
		},
	})
	require.NoError(t, err)

	lc := LastCommit(t)
	assert.Equal(t, "uplift,uplift@test.com,test commit", lc)
}

// LastCommit returns the latest log in the following format: author,email,message
func LastCommit(t *testing.T) string {
	t.Helper()

	out, err := git.Clean(git.Run("log", "-1", `--pretty=format:'%an,%ae,%B'`))
	require.NoError(t, err)

	return out
}

func trackFile(t *testing.T, name string) {
	err := ioutil.WriteFile(name, []byte(`hello, world`), 0644)
	require.NoError(t, err)

	err = git.Stage(name)
	require.NoError(t, err)
}

func TestRun_NoStagedFiles(t *testing.T) {
	git.InitRepo(t)

	err := Task{}.Run(&context.Context{
		CommitDetails: git.CommitDetails{
			Author:  "uplift",
			Email:   "uplift@test.com",
			Message: "test commit",
		},
	})
	require.NoError(t, err)

	lc := LastCommit(t)
	assert.NotContains(t, "test commit", lc)
}

func TestRun_NoGitRepository(t *testing.T) {
	git.MkTmpDir(t)

	err := Task{}.Run(&context.Context{})
	require.Error(t, err)
}
