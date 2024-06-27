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

func TestString(t *testing.T) {
	assert.Equal(t, "bumping files", Task{}.String())
}

func TestRun_NoBumpConfig(t *testing.T) {
	err := Task{}.Run(&context.Context{})
	assert.NoError(t, err)
}

func TestRun_NoStage(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, "test.txt", "version: 0.1.0")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: "test.txt",
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

	actual := ReadFile(t, "test.txt")
	assert.Equal(t, "version: 0.1.1", actual)

	status := gittest.PorcelainStatus(t)
	assert.ElementsMatch(t, []string{"?? test.txt"}, status)
}

func TestRun_GlobSupport(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, "a/1.json", `{"version": "0.1.0"}`)
	gittest.TempFile(t, "b/2.json", `{"version": "0.1.0"}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: "**/*.json",
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

	actual1 := ReadFile(t, "a/1.json")
	assert.Equal(t, `{"version": "0.1.1"}`, actual1)

	actual2 := ReadFile(t, "b/2.json")
	assert.Equal(t, `{"version": "0.1.1"}`, actual2)
}
