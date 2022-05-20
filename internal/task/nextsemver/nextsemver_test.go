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

package nextsemver

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
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
			git.Tag(tt.curVer)
			git.EmptyCommit(t, tt.commit)

			ctx := &context.Context{
				Prerelease: tt.prerelease,
				Metadata:   tt.metadata,
			}
			err := Task{}.Run(ctx)

			require.NoError(t, err)
			require.Equal(t, tt.expected, ctx.NextVersion.Raw)
		})
	}
}

func TestRun_NoSemanticBump(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "docs: started writing docs")

	ctx := &context.Context{}
	err := Task{}.Run(ctx)

	require.NoError(t, err)
	require.True(t, ctx.NoVersionChanged)
}

// TODO: test that prefix can be stripped from first version only

func TestRun_FirstVersion(t *testing.T) {
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
			expected: "0.0.1",
		},
		{
			name:     "MinorIncrement",
			commit:   "feat: a new feature",
			expected: "0.1.0",
		},
		{
			name:     "MajorIncrement",
			commit:   "feat!: a breaking change",
			expected: "1.0.0",
		},
		{
			name:       "MinorIncrementWithPrerelease",
			commit:     "feat: a new feature",
			prerelease: "beta.1",
			metadata:   "12345",
			expected:   "0.1.0-beta.1+12345",
		},
		{
			name: "BreakingChangeFooter",
			commit: `refactor: changed the cli
BREAKING CHANGE: no backwards compatibility support`,
			expected: "1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git.InitRepo(t)
			git.EmptyCommit(t, tt.commit)

			ctx := &context.Context{
				Prerelease: tt.prerelease,
				Metadata:   tt.metadata,
			}
			err := Task{}.Run(ctx)

			require.NoError(t, err)
			require.Equal(t, tt.expected, ctx.NextVersion.Raw)
		})
	}
}
