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

package semver

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
)

func TestDetectIncrement_BreakingFooter(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "docs: document about new breaking change",
		},
		{
			Message: "fix: annoying bug has now been fixed",
		},
		{
			Message: `refactor: changed a really important part of the app

BREAKING CHANGE: the cli has been completely refactored with no backwards compatibility`,
		},
		{
			Message: "docs(config): document new configuration option`",
		},
	}

	inc := DetectIncrement(log)
	assert.Equal(t, MajorIncrement, inc)
}

func TestDetectIncrement_BreakingFooterHyphenated(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "docs: document about new breaking change",
		},
		{
			Message: "fix: annoying bug has now been fixed",
		},
		{
			Message: `refactor: changed a really important part of the app

BREAKING-CHANGE: the cli has been completely refactored with no backwards compatibility`,
		},
		{
			Message: "docs(config): document new configuration option`",
		},
	}

	inc := DetectIncrement(log)
	assert.Equal(t, MajorIncrement, inc)
}

func TestDetectIncrement_BreakingBang(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "feat: a new snazzy feature has been added",
		},
		{
			Message: "fix: annoying bug has now been fixed",
		},
		{
			Message: "feat!: changed a really important part of the app",
		},
	}

	inc := DetectIncrement(log)
	assert.Equal(t, MajorIncrement, inc)
}

func TestDetectIncrement_Minor(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "ci: change to the existing workflow",
		},
		{
			Message: "fix: annoying bug has now been fixed",
		},
		{
			Message: "feat: shiny new feature has been added",
		},
	}

	inc := DetectIncrement(log)
	assert.Equal(t, MinorIncrement, inc)
}

func TestDetectIncrement_Patch(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "ci: change to the existing workflow",
		},
		{
			Message: "docs: updated documented to talk about fix",
		},
		{
			Message: "fix: small bug fixed",
		},
	}

	inc := DetectIncrement(log)
	assert.Equal(t, PatchIncrement, inc)
}

func TestDetectIncrement_NoIncrement(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "docs(ci): documented additional CI support",
		},
		{
			Message: "ci: sped up the existing build",
		},
		{
			Message: "docs(config): documented new configuration option",
		},
	}

	inc := DetectIncrement(log)
	assert.Equal(t, NoIncrement, inc)
}

func TestDetectIncrement_ChoreWithChangelog(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: `chore(deps): bumped version of vite to ^4.4.2

### Release Notes

<details>
<summary>vitejs/vite (vite)</summary>

##### Features

- feat: preview mode add keyboard shortcuts
- feat: emit event to handle chunk load errors
- feat!: update esbuild to 0.18.2`,
		},
	}
	inc := DetectIncrement(log)
	assert.Equal(t, NoIncrement, inc)
}
