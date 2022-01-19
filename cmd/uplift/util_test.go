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

	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/require"
)

func untaggedRepo(t *testing.T) {
	t.Helper()

	git.InitRepo(t)
	git.EmptyCommit(t, "feat: a new feature")
	require.Len(t, git.AllTags(), 0)
}

func taggedRepo(t *testing.T) {
	t.Helper()

	git.InitRepo(t)
	git.EmptyCommitAndTag(t, "1.0.0", "feat: a new feature")
}

func upliftConfigFile(t *testing.T, name string) {
	t.Helper()

	yml := "firstVersion: 1.0.0"

	err := ioutil.WriteFile(name, []byte(yml), 0644)
	require.NoError(t, err)
}
