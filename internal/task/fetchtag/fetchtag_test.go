package fetchtag

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "fetching all tags", Task{}.String())
}

func TestSkip(t *testing.T) {
	assert.True(t, Task{}.Skip(&context.Context{
		FetchTags: false,
	}))
}

func TestRun(t *testing.T) {
	log := `(tag: 0.2.1) fix: this is a fix
(tag: 0.2.0) feat: this was another feature
(tag: 0.1.0) feat: this was a feature`
	gittest.InitRepository(t, gittest.WithRemoteLog(log))
	require.Empty(t, gittest.Tags(t))

	err := Task{}.Run(&context.Context{})
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"0.1.0", "0.2.0", "0.2.1"}, gittest.Tags(t))
}
