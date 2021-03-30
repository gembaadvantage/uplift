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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchesMajorIncrementForType(t *testing.T) {
	tests := []struct {
		name   string
		commit string
		inc    increment
	}{
		{
			name:   "BuildBang",
			commit: "build!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "ChoreBang",
			commit: "chore!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "CIBang",
			commit: "ci!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "DocsBangPrefix",
			commit: "docs!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "FeatBang",
			commit: "feat!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "FixBang",
			commit: "fix!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "PerfBang",
			commit: "perf!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "RefactorBang",
			commit: "refactor!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "RevertBang",
			commit: "revert!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "StyleBang",
			commit: "style!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
		{
			name:   "TestBang",
			commit: "test!: Lorem ipsum dolor sit amet",
			inc:    majorIncrement,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inc := checkIncrement(tt.commit)
			if inc != tt.inc {
				t.Errorf("Expected %s but received %s", tt.inc, inc)
			}
		})
	}
}

func TestMatchesMinorIncrement(t *testing.T) {
	inc := checkIncrement("feat: Lorem ipsum dolor sit amet")
	assert.Equal(t, minorIncrement, inc)
}

func TestMatchesPatchIncrement(t *testing.T) {
	inc := checkIncrement("fix: Lorem ipsum dolor sit amet")
	assert.Equal(t, patchIncrement, inc)
}

func TestMatchesMajorIncrementForBreakingChangeFooter(t *testing.T) {
	inc := checkIncrement(`fix: Lorem ipsum dolor sit amet

BREAKING CHANGE: Lorem ipsum dolor sit amet`)
	assert.Equal(t, majorIncrement, inc)
}

func TestDoesNotMatchUnrecognisedType(t *testing.T) {
	inc := checkIncrement("feet: Lorem ipsum dolor sit amet")
	assert.Equal(t, noIncrement, inc)
}
