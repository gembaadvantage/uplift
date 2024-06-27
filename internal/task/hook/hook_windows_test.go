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
		"sed --posix -i 's/Doe/Smith/g' out.txt",
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
	LAST_COMMIT=$(git log -1 --pretty=format:'%B')
	echo -n $LAST_COMMIT > out.txt`
	os.Mkdir("subfolder", 0o755)
	os.WriteFile("subfolder/last-commit.sh", []byte(sh), 0o755)

	err := Exec(context.Background(), []string{"bash subfolder//last-commit.sh"}, ExecOptions{})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, gittest.InitialCommit, string(data))
}
