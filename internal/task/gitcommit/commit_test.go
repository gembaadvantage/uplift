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

	c, _ := git.LatestCommit()
	assert.Equal(t, "uplift", c.Author)
	assert.Equal(t, "uplift@test.com", c.Email)
	assert.Equal(t, "test commit", c.Message)
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

	c, _ := git.LatestCommit()
	assert.Equal(t, git.InitCommit, c.Message)
}

func TestRun_NoGitRepository(t *testing.T) {
	git.MkTmpDir(t)

	err := Task{}.Run(&context.Context{})
	require.Error(t, err)
}
