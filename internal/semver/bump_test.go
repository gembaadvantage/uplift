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

func TestIdentifyVersion(t *testing.T) {
	git.InitRepo(t)

	_, err := git.Tag("1.2.3")
	require.NoError(t, err)

	b := NewBumper(io.Discard, BumpOptions{FirstVersion: "0.1.0"})
	v, err := b.identifyVersion()

	require.NoError(t, err)
	assert.Equal(t, "1.2.3", v)
}

func TestIdentifyVersionFirstVersion(t *testing.T) {
	git.InitRepo(t)

	b := NewBumper(io.Discard, BumpOptions{FirstVersion: "0.1.0"})
	v, err := b.identifyVersion()

	require.NoError(t, err)
	assert.Equal(t, "0.1.0", v)
}

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		inc      increment
		expected string
	}{
		{
			name:     "MajorIncrement",
			version:  "1.2.3",
			inc:      majorIncrement,
			expected: "2.0.0",
		},
		{
			name:     "MinorIncrement",
			version:  "1.2.3",
			inc:      minorIncrement,
			expected: "1.3.0",
		},
		{
			name:     "PatchIncrement",
			version:  "1.2.3",
			inc:      patchIncrement,
			expected: "1.2.4",
		},
		{
			name:     "NoChange",
			version:  "1.2.3",
			inc:      noIncrement,
			expected: "1.2.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBumper(io.Discard, BumpOptions{})
			v, err := b.bumpVersion(tt.version, tt.inc)
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}

			if v != tt.expected {
				t.Errorf("Expected %s but received %s", tt.expected, v)
			}
		})
	}
}

func TestBumpVersionKeepsVPrefix(t *testing.T) {
	b := NewBumper(io.Discard, BumpOptions{})
	v, err := b.bumpVersion("v1.0.0", majorIncrement)

	require.NoError(t, err)
	assert.Equal(t, "v2.0.0", v)
}

func TestBumpVersionInvalidVersion(t *testing.T) {
	b := NewBumper(io.Discard, BumpOptions{})
	_, err := b.bumpVersion("1.0.B", minorIncrement)

	require.Error(t, err)
}
