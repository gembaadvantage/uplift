/*
Copyright (c) 2022 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
