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

func TestValidateBumpFilePathEmpty(t *testing.T) {
	cfg := Uplift{
		Bumps: []Bump{
			{
				File: "",
			},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Bumps[0].File")
}

func TestValidateBumpFilePathDoesNotResolve(t *testing.T) {
	cfg := Uplift{
		Bumps: []Bump{
			{
				File: "does-not-exist.txt",
			},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Bumps[0].File")
}

func TestValidateBumpFileRegexAndJSONEmpty(t *testing.T) {
	path := WriteFile(t, "test.txt")
	cfg := Uplift{
		Bumps: []Bump{
			{
				File: path,
			},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Bumps[0].Regex")
	require.ErrorContains(t, err, "Uplift.Bumps[0].JSON")
}

func TestValidateRegexBumpPatternEmpty(t *testing.T) {
	path := WriteFile(t, "test.txt")
	cfg := Uplift{
		Bumps: []Bump{
			{
				File: path,
				Regex: []RegexBump{
					{
						Pattern: "",
					},
				},
			},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Bumps[0].Regex[0].Pattern")
}

func TestValidateRegexBumpCountLessThanZero(t *testing.T) {
	path := WriteFile(t, "test.txt")
	cfg := Uplift{
		Bumps: []Bump{
			{
				File: path,
				Regex: []RegexBump{
					{
						Pattern: "version: $VERSION",
						Count:   -1,
					},
				},
			},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Bumps[0].Regex[0].Count")
}

func TestValidateJsonBumpPathEmpty(t *testing.T) {
	path := WriteFile(t, "test.txt")
	cfg := Uplift{
		Bumps: []Bump{
			{
				File: path,
				JSON: []JSONBump{
					{
						Path: "",
					},
				},
			},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Bumps[0].JSON[0].Path")
}

func TestValidateCommitAuthorNameAndEmailEmpty(t *testing.T) {
	cfg := Uplift{
		CommitAuthor: &CommitAuthor{},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.CommitAuthor.Name")
	require.ErrorContains(t, err, "Uplift.CommitAuthor.Email")
}

func TestValidateCommitAuthorNameEmpty(t *testing.T) {
	cfg := Uplift{
		CommitAuthor: &CommitAuthor{
			Name: "",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.CommitAuthor.Name")
}

func TestValidateCommitAuthorEmailEmpty(t *testing.T) {
	cfg := Uplift{
		CommitAuthor: &CommitAuthor{
			Name: "",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.CommitAuthor.Name")
}

func TestValidateCommitEmail(t *testing.T) {
	cfg := Uplift{
		CommitAuthor: &CommitAuthor{
			Email: "not-a-valid-email",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.CommitAuthor.Name")
}

func TestValidateChangelogSort(t *testing.T) {
	tests := []struct {
		name      string
		direction string
	}{
		{
			name:      "AscendingLowercase",
			direction: "asc",
		},
		{
			name:      "AscendingUppercase",
			direction: "ASC",
		},
		{
			name:      "DescendingLowercase",
			direction: "desc",
		},
		{
			name:      "DescendingUppercase",
			direction: "DESC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Uplift{
				Changelog: &Changelog{
					Sort: tt.direction,
				},
			}

			err := cfg.Validate()
			require.NoError(t, err)
		})
	}
}

func TestValidateChangelogSortUnsupported(t *testing.T) {
	cfg := Uplift{
		Changelog: &Changelog{
			Sort: "Bubble",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Changelog.Sort")
}

func TestValidateChangeLogExcludeEmpty(t *testing.T) {
	cfg := Uplift{
		Changelog: &Changelog{
			Exclude: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Changelog.Exclude[0]")
}

func TestValidateChangelogIncludeEmpty(t *testing.T) {
	cfg := Uplift{
		Changelog: &Changelog{
			Include: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Changelog.Include[0]")
}

func TestValidateChangelogAllEmpty(t *testing.T) {
	cfg := Uplift{
		Changelog: &Changelog{},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Changelog.Sort")
	require.ErrorContains(t, err, "Uplift.Changelog.Include")
	require.ErrorContains(t, err, "Uplift.Changelog.Exclude")
}

func TestValidateGitPushOptionEmpty(t *testing.T) {
	cfg := Uplift{
		Git: &Git{
			PushOptions: []GitPushOption{
				{
					Option: "",
				},
			},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Git.PushOptions[0].Option")
}

func TestValidateGitHubNonCompliantURL(t *testing.T) {
	cfg := Uplift{
		GitHub: &GitHub{
			URL: "not-a-url",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.GitHub.URL")
}

func TestValidateGitLabNonCompliantURL(t *testing.T) {
	cfg := Uplift{
		GitLab: &GitLab{
			URL: "not-a-url",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.GitLab.URL")
}

func TestValidateGiteaNonCompliantURL(t *testing.T) {
	cfg := Uplift{
		Gitea: &Gitea{
			URL: "not-a-url",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Gitea.URL")
}

func TestValidateHooksBeforeEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			Before: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Hooks.Before[0]")
}

func TestValidateHooksBeforeBumpEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			BeforeBump: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Hooks.BeforeBump[0]")
}

func TestValidateHooksBeforeTagEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			BeforeTag: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Hooks.BeforeTag[0]")
}

func TestValidateHooksBeforeChangelogEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			BeforeChangelog: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Hooks.BeforeChangelog[0]")
}

func TestValidateHooksAfterEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			After: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Hooks.After[0]")
}

func TestValidateHooksAfterBumpEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			AfterBump: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Hooks.AfterBump[0]")
}

func TestValidateHooksAfterTagEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			AfterTag: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Hooks.AfterTag[0]")
}

func TestValidateHooksAfterChangelogEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			AfterChangelog: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "Uplift.Hooks.AfterChangelog[0]")
}
