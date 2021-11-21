/*
Copyright (c) 2021 Gemba Advantage

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

package gittag

import (
	"fmt"
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tag := "1.1.0"
	git.InitRepo(t)
	git.EmptyCommitAndTag(t, "1.0.0", "commit")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: tag,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	ltag := git.LatestTag()
	assert.Equal(t, tag, ltag)
}

func TestRun_DryRunMode(t *testing.T) {
	tag := "v1.2.3"
	git.InitRepo(t)
	git.EmptyCommit(t, "commit")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: tag,
		},
		DryRun: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	ltag := git.LatestTag()
	assert.Empty(t, ltag)
}

func TestRun_NoVersionChange(t *testing.T) {
	tag := "1.0.0"
	git.InitRepo(t)
	git.EmptyCommitAndTag(t, tag, "commit")

	ctx := &context.Context{
		CurrentVersion: semver.Version{
			Raw: tag,
		},
		NextVersion: semver.Version{
			Raw: tag,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	ltag := git.LatestTag()
	assert.Equal(t, tag, ltag)
}

func TestRun_NoGitRepository(t *testing.T) {
	git.MkTmpDir(t)

	err := Task{}.Run(&context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
	})
	require.Error(t, err)
}

func TestRun_AnnotatedTag(t *testing.T) {
	tag := "1.1.0"
	git.InitRepo(t)
	git.EmptyCommitAndTag(t, "1.0.0", "commit")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: tag,
		},
		CommitDetails: git.CommitDetails{
			Author:  "joe.bloggs",
			Email:   "joe.bloggs@example.com",
			Message: "custom message",
		},
		Config: config.Uplift{
			AnnotatedTags: true,
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	out, _ := git.Clean(git.Run("for-each-ref", fmt.Sprintf("refs/tags/%s", tag),
		"--format='%(taggername):%(taggeremail):%(contents)'"))

	assert.Contains(t, out, fmt.Sprintf("%s:<%s>:%s",
		ctx.CommitDetails.Author, ctx.CommitDetails.Email, ctx.CommitDetails.Message))
}
