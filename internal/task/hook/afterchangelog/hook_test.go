package afterchangelog

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "after generating changelog", Task{}.String())
}

func TestSkip(t *testing.T) {
	cmd := []string{"echo 'HELLO'"}

	assert.True(t, Task{}.Skip(&context.Context{
		Config: config.Uplift{
			Hooks: &config.Hooks{
				Before:          cmd,
				BeforeBump:      cmd,
				BeforeTag:       cmd,
				BeforeChangelog: cmd,
				After:           cmd,
				AfterBump:       cmd,
				AfterTag:        cmd,
				AfterChangelog:  []string{},
			},
		},
	}))
}

func TestSkip_NoVersionChanged(t *testing.T) {
	assert.True(t, Task{}.Skip(&context.Context{
		NoVersionChanged: true,
	}))
}

func TestSkip_SkipChangelog(t *testing.T) {
	assert.True(t, Task{}.Skip(&context.Context{
		SkipChangelog: true,
	}))
}

func TestRun(t *testing.T) {
	// git.MkTmpDir(t)
	gittest.InitRepository(t)

	err := Task{}.Run(&context.Context{
		Config: config.Uplift{
			Hooks: &config.Hooks{
				AfterChangelog: []string{"touch a.out"},
			},
		},
	})
	require.NoError(t, err)
	assert.FileExists(t, "a.out")
}
