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

package tasks

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		repoTag  string
		expected string
	}{
		{
			name:     "RepositoryTag",
			repoTag:  "1.1.1",
			expected: "1.1.1",
		},
		{
			name:     "PrefixedRepositoryTag",
			repoTag:  "v0.2.1",
			expected: "v0.2.1",
		},
		{
			name:     "RepositoryNoTag",
			repoTag:  "",
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git.InitRepo(t)
			if tt.repoTag != "" {
				git.EmptyCommitAndTag(t, tt.repoTag, "testing")
			}

			ctx := &context.Context{}
			err := CurrentVersion{}.Run(ctx)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if ctx.CurrentVersion.Raw != tt.expected {
				t.Errorf("expected tag %s but received tag %s", tt.expected, ctx.CurrentVersion.Raw)
			}
		})
	}
}
