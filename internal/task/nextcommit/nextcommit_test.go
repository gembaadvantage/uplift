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

package nextcommit

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_DefaultCommitMessage(t *testing.T) {
	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, "ci(uplift): uplifted for version 0.1.0", ctx.CommitDetails.Message)
}

func TestRun_ImpersonatesAuthor(t *testing.T) {
	cd := git.CommitDetails{
		Author: "joe.bloggs",
		Email:  "joe.bloggs@example.com",
	}

	ctx := &context.Context{
		CommitDetails: cd,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, cd.Author, ctx.CommitDetails.Author)
	assert.Equal(t, cd.Email, ctx.CommitDetails.Email)
}

func TestRun_CustomCommit(t *testing.T) {
	ctx := &context.Context{
		Config: config.Uplift{
			CommitMessage: "ci(release): this is a custom message",
			CommitAuthor: config.CommitAuthor{
				Name:  "releasebot",
				Email: "releasebot@example.com",
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, "releasebot", ctx.CommitDetails.Author)
	assert.Equal(t, "releasebot@example.com", ctx.CommitDetails.Email)
	assert.Equal(t, "ci(release): this is a custom message", ctx.CommitDetails.Message)
}
