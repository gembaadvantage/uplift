/*
Copyright (c) 2021 Gemba Advantage

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
	rootCmd, err := newRootCmd([]string{}, os.Stdout)
	require.NoError(t, err)

	rootCmd.cmd.SetArgs([]string{"--dry-run"})
	err = rootCmd.cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, true, rootCmd.ctx.DryRun)
}

func TestRoot_DebugFlag(t *testing.T) {
	rootCmd, err := newRootCmd([]string{}, os.Stdout)
	require.NoError(t, err)

	rootCmd.cmd.SetArgs([]string{"--debug"})
	err = rootCmd.cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, true, rootCmd.ctx.Debug)
}

func TestRoot_NoPushFlag(t *testing.T) {
	rootCmd, err := newRootCmd([]string{}, os.Stdout)
	require.NoError(t, err)

	rootCmd.cmd.SetArgs([]string{"--no-push"})
	err = rootCmd.cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, true, rootCmd.ctx.NoPush)
}
