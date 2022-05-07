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

package beforehook

import (
	ctx "context"
	"io/ioutil"
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "running before hooks", Task{}.String())
}

func TestSkip(t *testing.T) {
	noHooks := []string{}

	assert.True(t, Task{}.Skip(&context.Context{
		Config: config.Uplift{
			Hooks: config.Hooks{
				Before: noHooks,
			},
		},
	}))
}

func TestRun_ShellCommands(t *testing.T) {
	git.MkTmpDir(t)

	tctx := &context.Context{
		Context: ctx.Background(),
		Config: config.Uplift{
			Hooks: config.Hooks{
				Before: []string{
					"echo -n 'JohnDoe' > out.txt",
					"sed -i '' 's/Doe/Smith/g' out.txt",
				},
			},
		},
	}

	err := Task{}.Run(tctx)
	require.NoError(t, err)

	data, err := ioutil.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "JohnSmith", string(data))
}

func TestRun_ShellScripts(t *testing.T) {
	git.InitRepo(t)

	// Generate a shell script
	sh := `#!/bin/bash
git checkout -b $BRANCH
CURRENT=$(git branch --show-current)
echo -n $CURRENT > out.txt`
	ioutil.WriteFile("switch-branch.sh", []byte(sh), 0755)

	tctx := &context.Context{
		Context: ctx.Background(),
		Config: config.Uplift{
			Hooks: config.Hooks{
				Before: []string{
					"BRANCH=testing ./switch-branch.sh",
				},
			},
		},
	}

	err := Task{}.Run(tctx)
	require.NoError(t, err)

	data, err := ioutil.ReadFile("out.txt")
	require.NoError(t, err)

	assert.Equal(t, "testing", string(data))
}
