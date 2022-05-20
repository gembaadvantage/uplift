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

package lastcommit

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, "scanning for conventional commit", Task{}.String())
}

func TestSkip(t *testing.T) {
	assert.False(t, Task{}.Skip(&context.Context{}))
}

// func TestRun(t *testing.T) {
// 	git.InitRepo(t)
// 	git.EmptyCommits(t, "feat: brand new feature", "Merge branch 'main' of https://github.com/org/repo")

// 	ctx := &context.Context{
// 		CurrentVersion: semver.Version{
// 			Raw: "",
// 		},
// 	}
// 	err := Task{}.Run(ctx)

// 	require.NoError(t, err)
// 	assert.Equal(t, "uplift", ctx.CommitDetails.Author)
// 	assert.Equal(t, "uplift@test.com", ctx.CommitDetails.Email)
// 	assert.Equal(t, "feat: brand new feature", ctx.CommitDetails.Message)
// }

// func TestRun_FromTag(t *testing.T) {
// 	c := `feat: brand new feature

// with some additional commit details`

// 	git.InitRepo(t)
// 	git.EmptyCommitsAndTag(t, "1.0.0", "feat: first new feature")
// 	git.EmptyCommits(t, c, "Merge branch 'main' of https://github.com/org/repo")

// 	ctx := &context.Context{
// 		CurrentVersion: semver.Version{
// 			Raw: "1.0.0",
// 		},
// 	}
// 	err := Task{}.Run(ctx)

// 	require.NoError(t, err)
// 	assert.Equal(t, "uplift", ctx.CommitDetails.Author)
// 	assert.Equal(t, "uplift@test.com", ctx.CommitDetails.Email)
// 	assert.Equal(t, c, ctx.CommitDetails.Message)
// }

// func TestRun_NoConventionalCommits(t *testing.T) {
// 	git.InitRepo(t)
// 	git.EmptyCommits(t, "first commit", "second commit", "third commit")

// 	ctx := &context.Context{
// 		CurrentVersion: semver.Version{
// 			Raw: "",
// 		},
// 	}
// 	err := Task{}.Run(ctx)

// 	require.NoError(t, err)
// 	assert.Equal(t, "uplift", ctx.CommitDetails.Author)
// 	assert.Equal(t, "uplift@test.com", ctx.CommitDetails.Email)
// 	assert.Equal(t, "third commit", ctx.CommitDetails.Message)
// }

// func TestRun_NoGitRepository(t *testing.T) {
// 	git.MkTmpDir(t)

// 	err := Task{}.Run(&context.Context{})
// 	require.Error(t, err)
// }
