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

func TestString(t *testing.T) {
	assert.Equal(t, "next semantic version", Task{}.String())
}

func TestSkip(t *testing.T) {
	assert.False(t, Task{}.Skip(&context.Context{}))
}

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		commit     string
		curVer     string
		prerelease string
		metadata   string
		expected   string
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
		{
			name:       "MinorIncrementWithPrerelease",
			commit:     "feat: a new feature",
			curVer:     "v0.1.0",
			prerelease: "beta.1",
			metadata:   "12345",
			expected:   "v0.2.0-beta.1+12345",
		},
		{
			name: "BreakingChangeFooter",
			commit: `refactor: changed the cli
BREAKING CHANGE: no backwards compatibility support`,
			curVer:   "v0.9.2",
			expected: "v1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git.InitRepo(t)
			if tt.commit != "" {
				git.EmptyCommitAndTag(t, tt.curVer, tt.commit)
			}

			ctx := &context.Context{
				CommitDetails: git.CommitDetails{
					Message: tt.commit,
				},
				CurrentVersion: semver.Version{
					Raw: tt.curVer,
				},
				Prerelease: tt.prerelease,
				Metadata:   tt.metadata,
			}
			err := Task{}.Run(ctx)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if ctx.NextVersion.Raw != tt.expected {
				t.Errorf("expected version %s but received version %s", tt.expected, ctx.NextVersion.Raw)
			}
		})
	}
}

func TestRun_FirstTagDefault(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: a new feature")

	ctx := &context.Context{
		CommitDetails: git.CommitDetails{
			Message: "feat: a new feature",
		},
	}
	err := Task{}.Run(ctx)

	require.NoError(t, err)
	assert.Equal(t, "0.1.0", ctx.NextVersion.Raw)
}

func TestRun_FirstTagDefaultFromConfig(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: a new feature")

	ctx := &context.Context{
		CommitDetails: git.CommitDetails{
			Message: "feat: a new feature",
		},
		Config: config.Uplift{
			FirstVersion: "1.0.0",
		},
	}
	err := Task{}.Run(ctx)

	require.NoError(t, err)
	assert.Equal(t, ctx.Config.FirstVersion, ctx.NextVersion.Raw)
}

func TestRun_FirstTagPrerelease(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: a new feature")

	ctx := &context.Context{
		CommitDetails: git.CommitDetails{
			Message: "feat: a new feature",
		},
		Config: config.Uplift{
			FirstVersion: "1.0.0",
		},
		Prerelease: "beta.1",
		Metadata:   "12345",
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0-beta.1+12345", ctx.NextVersion.Raw)
}
