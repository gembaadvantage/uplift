package gittag

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "tagging repository", Task{}.String())
}

func TestSkip(t *testing.T) {
	assert.True(t, Task{}.Skip(&context.Context{
		NoVersionChanged: true,
	}))
}

func TestRun(t *testing.T) {
	log := "(tag: 1.0.0) feat: an exciting new feature"
	gittest.InitRepository(t, gittest.WithLog(log))

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		NoPush: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	tags := gittest.RemoteTags(t)
	assert.ElementsMatch(t, []string{"1.0.0", "1.1.0"}, tags)
}

func TestRun_DryRunMode(t *testing.T) {
	gittest.InitRepository(t)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "v1.2.3",
		},
		DryRun: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	tags := gittest.RemoteTags(t)
	assert.Empty(t, tags)
}

func TestRun_NoVersionChange(t *testing.T) {
	log := "(tag: 0.1.0) feat: another feature"
	gittest.InitRepository(t, gittest.WithLog(log))

	ctx := &context.Context{
		CurrentVersion: semver.Version{
			Raw: "0.1.0",
		},
		NextVersion: semver.Version{
			Raw: "0.1.0",
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	tags := gittest.RemoteTags(t)
	assert.ElementsMatch(t, []string{"0.1.0"}, tags)
}

func TestRun_AnnotatedTag(t *testing.T) {
	log := "(tag: 0.2.0) feat: another feature"
	gittest.InitRepository(t, gittest.WithLog(log))

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.3.0",
		},
		CommitDetails: git.CommitDetails{
			Author: git.Person{
				Name:  "joe.bloggs",
				Email: "joe.bloggs@example.com",
			},
			Message: "custom message",
		},
		Config: config.Uplift{
			AnnotatedTags: true,
		},
		NoPush: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	out := gittest.MustExec(t, "git for-each-ref refs/tags/0.3.0 --format='%(taggername):%(taggeremail):%(contents)'")
	assert.Contains(t, out, fmt.Sprintf("%s:<%s>:%s", ctx.CommitDetails.Author.Name,
		ctx.CommitDetails.Author.Email, ctx.CommitDetails.Message))
}

func TestRun_PrintCurrentTag(t *testing.T) {
	gittest.InitRepository(t)

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		PrintCurrentTag: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0", buf.String())
	tags := gittest.RemoteTags(t)
	assert.Empty(t, tags)
}

func TestRun_PrintNextTag(t *testing.T) {
	gittest.InitRepository(t)

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
		PrintNextTag: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0", buf.String())
	tags := gittest.RemoteTags(t)
	assert.Empty(t, tags)
}

func TestRun_PrintCurrentAndNextTag(t *testing.T) {
	gittest.InitRepository(t)

	var buf bytes.Buffer
	ctx := &context.Context{
		Out: &buf,
		CurrentVersion: semver.Version{
			Raw: "1.0.0",
		},
		NextVersion: semver.Version{
			Raw: "1.1.0",
		},
		PrintCurrentTag: true,
		PrintNextTag:    true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	assert.Equal(t, "1.0.0 1.1.0", buf.String())
	tags := gittest.RemoteTags(t)
	assert.Empty(t, tags)
}

func TestFilterPushOptions(t *testing.T) {
	pushOpts := []config.GitPushOption{
		{
			Option: "option1",
		},
		{
			Option:  "option2",
			SkipTag: true,
		},
		{
			Option:     "option3",
			SkipBranch: true,
		},
	}

	filtered := filterPushOptions(pushOpts)
	assert.Len(t, filtered, 2)
	assert.Equal(t, []string{"option1", "option3"}, filtered)
}
