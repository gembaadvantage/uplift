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
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot_DryRunFlag(t *testing.T) {
	ctx := &context.Context{}
	cmd, err := newRootCmd([]string{}, ctx)
	require.NoError(t, err)

	cmd.SetArgs([]string{"--dry-run"})
	err = cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, true, ctx.DryRun)
}

func TestRoot_DebugFlag(t *testing.T) {
	ctx := &context.Context{}
	cmd, err := newRootCmd([]string{}, ctx)
	require.NoError(t, err)

	cmd.SetArgs([]string{"--debug"})
	err = cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, true, ctx.Debug)
}

func TestRoot_NoPushFlag(t *testing.T) {
	ctx := &context.Context{}
	cmd, err := newRootCmd([]string{}, ctx)
	require.NoError(t, err)

	cmd.SetArgs([]string{"--no-push"})
	err = cmd.Execute()
	require.NoError(t, err)

	assert.Equal(t, true, ctx.NoPush)
}
