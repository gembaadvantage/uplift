package bump

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_JSONNonMatchingPath(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, "test.json", `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: "test.json",
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
	gittest.InitRepository(t)
	gittest.TempFile(t, "example.json", `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "v0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: "example.json",
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

	actual := ReadFile(t, "example.json")
	assert.Equal(t, `{"version": "0.1.0"}`, actual)
}

func TestRun_JSONStrictSemVer(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, "example.json", `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "v0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: "example.json",
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

	actual := ReadFile(t, "example.json")
	assert.Equal(t, `{"version": "0.2.0"}`, actual)
}

func TestRun_JSONDryRun(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, "test.json", `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: "test.json",
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

	actual := ReadFile(t, "test.json")
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
	gittest.InitRepository(t)
	gittest.TempFile(t, "temp.json", `{
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
					File: "temp.json",
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

	actual := ReadFile(t, "temp.json")
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
