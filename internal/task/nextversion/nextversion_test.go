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

package nextversion

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		commit   string
		curVer   string
		expected string
	}{
		{
			name:     "PatchIncrement",
			commit:   "fix: a new fix",
			curVer:   "0.1.0",
			expected: "0.1.1",
		},
		{
			name:     "MinorIncrement",
			commit:   "feat: a new feature",
			curVer:   "v0.3.0",
			expected: "v0.4.0",
		},
		{
			name:     "MajorIncrement",
			commit:   "feat!: a breaking change",
			curVer:   "1.0.0",
			expected: "2.0.0",
		},
		{
			name:     "NoChange",
			commit:   "docs: change to readme",
			curVer:   "v0.1.0",
			expected: "v0.1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git.InitRepo(t)
			if tt.commit != "" {
				git.EmptyCommitAndTag(t, tt.curVer, tt.commit)
			}

			ctx := &context.Context{
				CurrentVersion: semver.Version{
					Raw: tt.curVer,
				},
			}
			err := Task{}.Run(ctx)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if ctx.NextVersion.Raw != tt.expected {
				t.Errorf("expected tag %s but received tag %s", tt.expected, ctx.NextVersion.Raw)
			}
		})
	}
}

func TestRun_FirstTagDefault(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: a new feature")

	ctx := &context.Context{}
	err := Task{}.Run(ctx)

	require.NoError(t, err)
	assert.Equal(t, "0.1.0", ctx.NextVersion.Raw)
}

func TestRun_FirstTagDefaultFromConfig(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: a new feature")

	ctx := &context.Context{
		Config: config.Uplift{
			FirstVersion: "1.0.0",
		},
	}
	err := Task{}.Run(ctx)

	require.NoError(t, err)
	assert.Equal(t, ctx.Config.FirstVersion, ctx.NextVersion.Raw)
}

func TestRun_NoGitRepository(t *testing.T) {
	git.MkTmpDir(t)

	err := Task{}.Run(&context.Context{})
	require.Error(t, err)
}
