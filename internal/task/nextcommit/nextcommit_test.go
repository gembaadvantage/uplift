package nextcommit

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
	assert.Equal(t, "building next commit", Task{}.String())
}

func TestSkip(t *testing.T) {
	assert.True(t, Task{}.Skip(&context.Context{
		NoVersionChanged: true,
	}))
}

func TestRun(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.name", "", "user.email", "")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, "uplift-bot", ctx.CommitDetails.Author.Name)
	assert.Equal(t, "uplift@gembaadvantage.com", ctx.CommitDetails.Author.Email)
	assert.Equal(t, "ci(uplift): uplifted for version 0.1.0", ctx.CommitDetails.Message)
}

func TestRun_GitAuthorConfig(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.name", "john.smith", "user.email", "john.smith@testing.com")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, "john.smith", ctx.CommitDetails.Author.Name)
	assert.Equal(t, "john.smith@testing.com", ctx.CommitDetails.Author.Email)
	assert.Equal(t, "ci(uplift): uplifted for version 0.1.0", ctx.CommitDetails.Message)
}

func TestRun_CustomCommitDetails(t *testing.T) {
	gittest.InitRepository(t)
	gittest.ConfigSet(t, "user.name", "", "user.email", "")

	ctx := &context.Context{
		Config: config.Uplift{
			CommitMessage: "ci(release): this is a custom message",
			CommitAuthor: &config.CommitAuthor{
				Name:  "releasebot",
				Email: "releasebot@example.com",
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, "releasebot", ctx.CommitDetails.Author.Name)
	assert.Equal(t, "releasebot@example.com", ctx.CommitDetails.Author.Email)
	assert.Equal(t, "ci(release): this is a custom message", ctx.CommitDetails.Message)
}

func TestRun_CustomCommitWithVersionToken(t *testing.T) {
	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			CommitMessage: "ci(release): a release for $VERSION",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, "ci(release): a release for 0.1.1", ctx.CommitDetails.Message)
}
