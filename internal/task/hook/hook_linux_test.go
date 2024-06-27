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
		"sed -i 's/Doe/Smith/g' out.txt",
	}

	err := Exec(context.Background(), cmds, ExecOptions{})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "JohnSmith", string(data))
}

func TestExec_ShellScripts(t *testing.T) {
	log := "(tag: 1.0.0) feat: first release"
	gittest.InitRepository(t, gittest.WithLog(log))

	// Generate a shell script
	sh := `#!/bin/bash
	LATEST_TAG=$(git for-each-ref "refs/tags/*.*.*" --sort=-v:creatordate --format='%(refname:short)')
	echo -n $LATEST_TAG > out.txt`
	os.WriteFile("latest-tag.sh", []byte(sh), 0o755)

	err := Exec(context.Background(), []string{"./latest-tag.sh"}, ExecOptions{})
	require.NoError(t, err)

	data, err := os.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "1.0.0", string(data))
}
