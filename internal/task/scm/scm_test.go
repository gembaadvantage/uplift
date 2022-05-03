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

package scm

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "detecting scm provider", Task{}.String())
}

func TestSkip(t *testing.T) {
	assert.False(t, Task{}.Skip(&context.Context{}))
}

func TestRun(t *testing.T) {
	tests := []struct {
		name      string
		remote    string
		provider  git.SCM
		tagURL    string
		commitURL string
	}{
		{
			name:      "GitHub",
			remote:    "https://github.com/owner/repository.git",
			provider:  git.GitHub,
			tagURL:    "https://github.com/owner/repository/releases/tag/{{.Ref}}",
			commitURL: "https://github.com/owner/repository/commit/{{.Hash}}",
		},
		{
			name:      "GitLab",
			provider:  git.GitLab,
			remote:    "https://gitlab.com/owner/repository.git",
			tagURL:    "https://gitlab.com/owner/repository/-/tags/{{.Ref}}",
			commitURL: "https://gitlab.com/owner/repository/-/commit/{{.Hash}}",
		},
		{
			name:      "CodeCommit",
			provider:  git.CodeCommit,
			remote:    "https://git-codecommit.eu-west-1.amazonaws.com/v1/repos/repository",
			tagURL:    "https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/repository/browse/refs/tags/{{.Ref}}?region=eu-west-1",
			commitURL: "https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/repository/commit/{{.Hash}}?region=eu-west-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git.InitRepo(t)
			git.RemoteOrigin(t, tt.remote)

			ctx := &context.Context{}
			err := Task{}.Run(ctx)

			require.NoError(t, err)
			require.Equal(t, ctx.SCM.Provider, tt.provider)
			require.Equal(t, ctx.SCM.TagURL, tt.tagURL)
			require.Equal(t, ctx.SCM.CommitURL, tt.commitURL)
		})
	}
}

func TestRun_GiteaSelfHosted(t *testing.T) {
	git.InitRepo(t)
	git.RemoteOrigin(t, "https://gitea.com/owner/repository.git")

	ctx := &context.Context{
		Config: config.Uplift{
			Gitea: config.Gitea{
				Host: "gitea.com",
			},
		},
	}
	err := Task{}.Run(ctx)

	require.NoError(t, err)
	assert.Equal(t, ctx.SCM.Provider, git.Gitea)
	assert.Equal(t, ctx.SCM.TagURL, "https://gitea.com/owner/repository/releases/tag/{{.Ref}}")
	assert.Equal(t, ctx.SCM.CommitURL, "https://gitea.com/owner/repository/commit/{{.Hash}}")
}

func TestRun_NoRemoteSet(t *testing.T) {
	git.MkTmpDir(t)

	err := Task{}.Run(&context.Context{})
	require.Error(t, err)
}

func TestRun_UnrecognisedSCM(t *testing.T) {
	git.InitRepo(t)
	git.RemoteOrigin(t, "https://unrecognised.com/owner/repository.git")

	ctx := &context.Context{}
	err := Task{}.Run(ctx)

	require.NoError(t, err)
	assert.Equal(t, git.Unrecognised, ctx.SCM.Provider)
}
