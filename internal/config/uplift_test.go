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

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("missing_file.yml")
	require.Error(t, err)
}

func TestLoadUnsupportedYaml(t *testing.T) {
	path := WriteFile(t, `
unrecognised_field: ""`)

	_, err := Load(path)
	require.Error(t, err)
}

func TestLoadInvalidYaml(t *testing.T) {
	path := WriteFile(t, `
doc: [`)

	_, err := Load(path)
	require.Error(t, err)
}

func WriteFile(t *testing.T, s string) string {
	t.Helper()

	current, err := os.Getwd()
	require.NoError(t, err)

	file, err := os.CreateTemp(current, "*")
	require.NoError(t, err)

	_, err = file.WriteString(s)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	t.Cleanup(func() {
		require.NoError(t, os.Remove(file.Name()))
	})

	return file.Name()
}

func TestUnmarshalGitPushOption(t *testing.T) {
	path := WriteFile(t, `
git:
  pushOptions:
    - custom-option
`)

	cfg, err := Load(path)

	require.NoError(t, err)
	assert.Len(t, cfg.Git.PushOptions, 1)

	opt := cfg.Git.PushOptions[0]
	assert.Equal(t, "custom-option", opt.Option)
	assert.False(t, opt.SkipBranch)
	assert.False(t, opt.SkipTag)
}

func TestUnmarshalGitPushOptionComplex(t *testing.T) {
	path := WriteFile(t, `
git:
  pushOptions:
    - option: custom-option-1
      skipTag: true
    - option: custom-option-2
      skipBranch: true
`)

	cfg, err := Load(path)

	require.NoError(t, err)
	assert.Len(t, cfg.Git.PushOptions, 2)

	opt1 := cfg.Git.PushOptions[0]
	assert.Equal(t, "custom-option-1", opt1.Option)
	assert.True(t, opt1.SkipTag)
	opt2 := cfg.Git.PushOptions[1]
	assert.Equal(t, "custom-option-2", opt2.Option)
	assert.True(t, opt2.SkipBranch)
}
