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

package semver

import (
	"io"
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBump(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		commit     string
		newVersion string
	}{
		{
			name:       "MajorIncrement",
			version:    "1.2.3",
			commit:     "refactor!: Lorem ipsum dolor sit amet",
			newVersion: "2.0.0",
		},
		{
			name:    "MajorIncrementBreakingChangeFooter",
			version: "1.2.3",
			commit: `refactor: Lorem ipsum dolor sit amet

BREAKING CHANGE: Lorem ipsum dolor sit amet`,
			newVersion: "2.0.0",
		},
		{
			name:       "MinorIncrement",
			version:    "1.2.3",
			commit:     "feat: Lorem ipsum dolor sit amet",
			newVersion: "1.3.0",
		},
		{
			name:       "PatchIncrement",
			version:    "1.2.3",
			commit:     "fix(db): Lorem ipsum dolor sit amet",
			newVersion: "1.2.4",
		},
		{
			name:       "NoChange",
			version:    "1.2.3",
			commit:     "chore: Lorem ipsum dolor sit amet",
			newVersion: "1.2.3",
		},
		{
			name:       "KeepsVersionPrefix",
			version:    "v1.1.1",
			commit:     "fix: Lorem ipsum dolor sit amet",
			newVersion: "v1.1.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git.InitRepo(t)
			tag(t, tt.version)

			git.EmptyCommit(t, tt.commit)

			b := NewBumper(io.Discard, BumpOptions{})
			err := b.Bump()
			require.NoError(t, err)

			v := git.LatestTag()

			if v != tt.newVersion {
				t.Errorf("Expected %s but received %s", tt.newVersion, v)
			}
		})
	}
}

func tag(t *testing.T, tag string) {
	t.Helper()

	err := git.Tag(tag)
	require.NoError(t, err)
}

func TestBumpInvalidVersion(t *testing.T) {
	git.InitRepo(t)
	tag(t, "1.0.B")
	git.EmptyCommit(t, "feat: Lorem ipsum dolor sit amet")

	b := NewBumper(io.Discard, BumpOptions{})
	err := b.Bump()

	require.Error(t, err)
}

func TestBumpFirstVersion(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: Lorem ipsum dolor sit amet")

	b := NewBumper(io.Discard, BumpOptions{FirstVersion: "0.1.0"})
	err := b.Bump()
	require.NoError(t, err)

	v := git.LatestTag()
	assert.Equal(t, "0.1.0", v)
}

func TestBumpEmptyRepo(t *testing.T) {
	git.InitRepo(t)

	b := NewBumper(io.Discard, BumpOptions{})
	err := b.Bump()

	require.NoError(t, err)
}

func TestBumpNotGitRepo(t *testing.T) {
	git.MkTmpDir(t)

	b := NewBumper(io.Discard, BumpOptions{})
	err := b.Bump()

	require.Error(t, err)
	assert.Error(t, err, "current directory must be a git repo")
}

func TestBumpAlwaysUseLatestCommit(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommits(t,
		"feat: Lorem ipsum dolor sit amet",
		"fix: Lorem ipsum dolor sit amet",
		"docs: Lorem ipsum dolor sit amet")

	b := NewBumper(io.Discard, BumpOptions{})
	err := b.Bump()

	require.NoError(t, err)
	assert.Equal(t, "", git.LatestTag())
}
