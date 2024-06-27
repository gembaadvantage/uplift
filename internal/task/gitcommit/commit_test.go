package gitcommit

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	git "github.com/purpleclay/gitz"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "committing changes", Task{}.String())
}

func TestSkip(t *testing.T) {
	tests := []struct {
		name string
		ctx  *context.Context
	}{
		{
			name: "DryRun",
			ctx: &context.Context{
				DryRun: true,
			},
		},
		{
			name: "NoVersionChanged",
			ctx: &context.Context{
				NoVersionChanged: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, Task{}.Skip(tt.ctx))
		})
	}
}

func TestRun(t *testing.T) {
	gittest.InitRepository(t, gittest.WithStagedFiles("test.txt"))

	err := Task{}.Run(&context.Context{
		CommitDetails: git.CommitDetails{
			Author: git.Person{
				Name:  "uplift",
				Email: "uplift@test.com",
			},
			Message: "test commit",
		},
	})
	require.NoError(t, err)

	lc := gittest.LastCommit(t)
	assert.Equal(t, "uplift", lc.AuthorName)
	assert.Equal(t, "uplift@test.com", lc.AuthorEmail)
	assert.Equal(t, "test commit", lc.Message)
}

func TestRun_NoStagedFiles(t *testing.T) {
	gittest.InitRepository(t)

	err := Task{}.Run(&context.Context{
		CommitDetails: git.CommitDetails{
			Author: git.Person{
				Name:  "uplift",
				Email: "uplift@test.com",
			},
			Message: "test commit",
		},
	})
	require.NoError(t, err)

	lc := gittest.LastCommit(t)
	assert.NotEqual(t, "test commit", lc.Message)
}

func TestFilterPushOptions(t *testing.T) {
	pushOpts := []config.GitPushOption{
		{
			Option: "option1",
		},
		{
			Option:     "option2",
			SkipBranch: true,
		},
		{
			Option:  "option3",
			SkipTag: true,
		},
	}

	filtered := filterPushOptions(pushOpts)
	assert.Len(t, filtered, 2)
	assert.Equal(t, []string{"option1", "option3"}, filtered)
}
