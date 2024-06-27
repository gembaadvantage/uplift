package semver

import (
	"testing"

	git "github.com/purpleclay/gitz"
	"github.com/stretchr/testify/assert"
)

func TestParseLog_BreakingFooter(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "docs: document about new breaking change",
		},
		{
			Message: "fix: annoying bug has now been fixed",
		},
		{
			Message: `refactor: changed a really important part of the app

BREAKING CHANGE: the cli has been completely refactored with no backwards compatibility`,
		},
		{
			Message: "docs(config): document new configuration option",
		},
	}

	inc := ParseLog(log)
	assert.Equal(t, MajorIncrement, inc)
}

func TestParseLog_BreakingBang(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "feat: a new snazzy feature has been added",
		},
		{
			Message: "fix: annoying bug has now been fixed",
		},
		{
			Message: "feat!: changed a really important part of the app",
		},
	}

	inc := ParseLog(log)
	assert.Equal(t, MajorIncrement, inc)
}

func TestParseLog_Minor(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "ci: change to the existing workflow",
		},
		{
			Message: "fix: annoying bug has now been fixed",
		},
		{
			Message: "feat: shiny new feature has been added",
		},
	}

	inc := ParseLog(log)
	assert.Equal(t, MinorIncrement, inc)
}

func TestParseLog_Patch(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "ci: change to the existing workflow",
		},
		{
			Message: "docs: updated documented to talk about fix",
		},
		{
			Message: "fix: small bug fixed",
		},
	}

	inc := ParseLog(log)
	assert.Equal(t, PatchIncrement, inc)
}

func TestParseLog_NoIncrement(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "docs(ci): documented additional CI support",
		},
		{
			Message: "ci: sped up the existing build",
		},
		{
			Message: "docs(config): documented new configuration option",
		},
	}

	inc := ParseLog(log)
	assert.Equal(t, NoIncrement, inc)
}

func TestParseLog_RenovateMultilineCommit(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: "fix: bug within search has been fixed",
		},
		{
			Message: `chore(deps): update dependency aws/aws-cdk to v2.90.0

## Changelog:

feat: a brand new feature has been added
feat!: a breaking change to the interface
refactor: tidy up some bits of the code`,
		},
	}

	inc := ParseLog(log)
	assert.Equal(t, PatchIncrement, inc)
}

func TestParseLog_TrimHeader(t *testing.T) {
	log := []git.LogEntry{
		{
			Message: `this line that should be ignored
this line that should also be ignored
feat: shiny new feature has been added`,
		},
	}

	inc := ParseLogWithOptions(log, ParseOptions{TrimHeader: true})
	assert.Equal(t, MinorIncrement, inc)
}
