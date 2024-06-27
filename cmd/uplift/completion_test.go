package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletion_Bash(t *testing.T) {
	var buf bytes.Buffer
	cmd := newCompletionCmd(&buf)
	cmd.SetArgs([]string{"bash"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "bash completion V2 for completion ")
}

func TestCompletion_Zsh(t *testing.T) {
	var buf bytes.Buffer
	cmd := newCompletionCmd(&buf)
	cmd.SetArgs([]string{"zsh"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "zsh completion for completion")
}

func TestCompletion_ZshNoDescriptions(t *testing.T) {
	var buf bytes.Buffer
	cmd := newCompletionCmd(&buf)
	cmd.SetArgs([]string{"zsh", "--no-descriptions"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "zsh completion for completion")
	assert.Contains(t, buf.String(), "__completeNoDesc")
}

func TestCompletion_Fish(t *testing.T) {
	var buf bytes.Buffer
	cmd := newCompletionCmd(&buf)
	cmd.SetArgs([]string{"fish"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "fish completion for completion")
}
