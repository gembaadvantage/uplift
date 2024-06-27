package hook

import (
	"context"
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec_ShellCommands(t *testing.T) {
	gittest.InitRepository(t)

	cmds := []string{
		"echo -n 'JohnDoe' > out.txt",
		"sed -i '' 's/Doe/Smith/g' out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "JohnSmith", string(data))
}

func TestExec_ShellScripts(t *testing.T) {
	gittest.InitRepository(t)

	// Generate a shell script
	sh := `#!/bin/bash
	git checkout -b $BRANCH
	CURRENT=$(git branch --show-current)
	echo -n $CURRENT > out.txt`
	os.WriteFile("switch-branch.sh", []byte(sh), 0o755)

	err := Exec(context.Background(), []string{"BRANCH=testing ./switch-branch.sh"}, ExecOptions{})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "testing", string(data))
}
