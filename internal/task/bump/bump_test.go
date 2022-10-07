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

package bump

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
	assert.Equal(t, "bumping files", Task{}.String())
}

func TestRun_NoBumpConfig(t *testing.T) {
	err := Task{}.Run(&context.Context{})
	assert.NoError(t, err)
}

func TestRun_NotGitRepository(t *testing.T) {
	git.MkTmpDir(t)
	file := WriteFile(t, "version: 0.1.0")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	assert.EqualError(t, err, "fatal: not a git repository (or any of the parent directories): .git")
}

func TestRun_NoStage(t *testing.T) {
	git.InitRepo(t)
	file := WriteFile(t, "version: 0.1.0")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
						},
					},
				},
			},
		},
		NoStage: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, file)
	assert.Equal(t, "version: 0.1.1", actual)

	staged, _ := git.Staged()
	assert.Empty(t, staged)
}
