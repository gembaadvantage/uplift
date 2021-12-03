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
	"io/ioutil"
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelease(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: this is a release")
	data := testFileWithConfig(t, "test.txt", ".uplift.yml")

	cfg, _ := config.Load(".uplift.yml")
	cmd := newReleaseCmd(&context.Context{Config: cfg})

	err := cmd.Execute()
	require.NoError(t, err)

	tags := git.AllTags()
	assert.Len(t, tags, 1)

	actual, err := ioutil.ReadFile("test.txt")
	require.NoError(t, err)
	assert.NotEqual(t, string(data), string(actual))
}

func TestRelease_CheckFlag(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: this is a release")

	cmd := newReleaseCmd(&context.Context{})
	cmd.SetArgs([]string{"--check"})

	err := cmd.Execute()
	require.NoError(t, err)
}

func TestRelease_CheckFlagNoRelease(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "ci: not a release")

	cmd := newReleaseCmd(&context.Context{})
	cmd.SetArgs([]string{"--check"})

	err := cmd.Execute()
	require.Error(t, err)
}
