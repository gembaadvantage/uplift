package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot_DryRunFlag(t *testing.T) {
	rootCmd := newRootCmd(os.Stdout)

	rootCmd.Cmd.SetArgs([]string{"--dry-run"})
	err := rootCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, rootCmd.Opts.DryRun)
}

func TestRoot_DebugFlag(t *testing.T) {
	rootCmd := newRootCmd(os.Stdout)

	rootCmd.Cmd.SetArgs([]string{"--debug"})
	err := rootCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, rootCmd.Opts.Debug)
}

func TestRoot_NoPushFlag(t *testing.T) {
	rootCmd := newRootCmd(os.Stdout)

	rootCmd.Cmd.SetArgs([]string{"--no-push"})
	err := rootCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, rootCmd.Opts.NoPush)
}

func TestRoot_ConfigDir(t *testing.T) {
	rootCmd := newRootCmd(os.Stdout)

	rootCmd.Cmd.SetArgs([]string{"--config-dir", "custom"})
	err := rootCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, "custom", rootCmd.Opts.ConfigDir)
}

func TestRoot_IgnoreDetachedFlag(t *testing.T) {
	rootCmd := newRootCmd(os.Stdout)

	rootCmd.Cmd.SetArgs([]string{"--ignore-detached"})
	err := rootCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, rootCmd.Opts.IgnoreDetached)
}

func TestRoot_IgnoreShallowFlag(t *testing.T) {
	rootCmd := newRootCmd(os.Stdout)

	rootCmd.Cmd.SetArgs([]string{"--ignore-shallow"})
	err := rootCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, rootCmd.Opts.IgnoreShallow)
}

func TestRoot_NoStage(t *testing.T) {
	rootCmd := newRootCmd(os.Stdout)

	rootCmd.Cmd.SetArgs([]string{"--no-stage"})
	err := rootCmd.Cmd.Execute()
	require.NoError(t, err)

	assert.True(t, rootCmd.Opts.NoStage)
}
