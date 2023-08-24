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

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const (
	HookDir             = "hooks/"
	BeforeFile          = HookDir + "before.out"
	BeforeBumpFile      = HookDir + "beforeBump.out"
	BeforeTagFile       = HookDir + "beforeTag.out"
	BeforeChangelogFile = HookDir + "beforeChangelog.out"
	AfterFile           = HookDir + "after.out"
	AfterBumpFile       = HookDir + "afterBump.out"
	AfterTagFile        = HookDir + "afterTag.out"
	AfterChangelogFile  = HookDir + "afterChangelog.out"
)

func noChangesPushed() *globalOptions {
	return &globalOptions{NoPush: true}
}

func numHooksExecuted(t *testing.T) int {
	t.Helper()

	de, err := os.ReadDir(HookDir)
	require.NoError(t, err)

	return len(de)
}

// Ensures all available hooks are configured. Each hook will create an empty
// file based on the defined test files. This should make verification
// of hooks easy, by checking the number of files touched and their respective
// filenames
func configWithHooks(t *testing.T) {
	t.Helper()

	err := os.Mkdir(HookDir, 0o755)
	require.NoError(t, err)

	cfg := &config.Uplift{
		Hooks: &config.Hooks{
			Before:          []string{"touch " + BeforeFile},
			BeforeBump:      []string{"touch " + BeforeBumpFile},
			BeforeTag:       []string{"touch " + BeforeTagFile},
			BeforeChangelog: []string{"touch " + BeforeChangelogFile},
			After:           []string{"touch " + AfterFile},
			AfterBump:       []string{"touch " + AfterBumpFile},
			AfterTag:        []string{"touch " + AfterTagFile},
			AfterChangelog:  []string{"touch " + AfterChangelogFile},
		},
	}
	data, err := yaml.Marshal(&cfg)
	require.NoError(t, err)

	err = os.WriteFile(".uplift.yml", data, 0o644)
	require.NoError(t, err)

	err = os.WriteFile(".gitignore", []byte(HookDir), 0o644)
	require.NoError(t, err)

	// Ensure files are committed to prevent dirty repository
	gittest.StageFile(t, ".gitignore")
	gittest.StageFile(t, ".uplift.yml")
	gittest.Commit(t, "chore: add files")
}
