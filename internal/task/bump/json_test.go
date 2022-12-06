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

func TestRun_JSONNonMatchingPath(t *testing.T) {
	path := WriteTempFile(t, `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: path,
					JSON: []config.JSONBump{
						{
							Path: "nomatch",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	assert.EqualError(t, err, "no version matched in file")
}

func TestRun_JSONNotAllPathsMatch(t *testing.T) {
	path := WriteTempFile(t, `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "v0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: path,
					JSON: []config.JSONBump{
						{
							Path:   "version",
							SemVer: true,
						},
						{
							Path:   "ver",
							SemVer: true,
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.Error(t, err)

	actual := ReadFile(t, path)
	assert.Equal(t, `{"version": "0.1.0"}`, actual)
}

func TestRun_JSONStrictSemVer(t *testing.T) {
	git.InitRepo(t)
	path := WriteTempFile(t, `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "v0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: path,
					JSON: []config.JSONBump{
						{
							Path:   "version",
							SemVer: true,
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, path)
	assert.Equal(t, `{"version": "0.2.0"}`, actual)
}

func TestRun_JSONDryRun(t *testing.T) {
	git.InitRepo(t)
	path := WriteTempFile(t, `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: path,
					JSON: []config.JSONBump{
						{
							Path: "version",
						},
					},
				},
			},
		},
		DryRun: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, path)
	assert.Equal(t, `{"version": "0.1.0"}`, actual)
}

func TestRun_JSONFileDoesNotExist(t *testing.T) {
	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: "missing.txt",
					JSON: []config.JSONBump{
						{
							Path: "version",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	assert.Error(t, err)
}

func TestRun_PackageJson(t *testing.T) {
	git.InitRepo(t)

	file := WriteTempFile(t, `{
  "name": "test",
  "version": "0.1.0",
  "bin": {
    "test": "bin/test.js"
  },
  "scripts": {
    "build": "tsc",
  },
  "devDependencies": {
    "typescript": "~3.7.2"
  },
  "dependencies": {}
}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
		CommitDetails: commit,
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					JSON: []config.JSONBump{
						{
							Path: "version",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, file)
	assert.Equal(t, `{
  "name": "test",
  "version": "1.0.0",
  "bin": {
    "test": "bin/test.js"
  },
  "scripts": {
    "build": "tsc",
  },
  "devDependencies": {
    "typescript": "~3.7.2"
  },
  "dependencies": {}
}`, actual)
}
