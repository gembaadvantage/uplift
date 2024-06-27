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

func TestUnmarshalGitPushOptionInvalid(t *testing.T) {
	path := WriteFile(t, `
git:
  pushOptions:
    - invalid: option
`)

	_, err := Load(path)

	require.Error(t, err)
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
	require.ErrorContains(t, err, "field 'Uplift.Bumps[0].File' contains a value that does not meet the minimum expected length of '1'")
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
	require.ErrorContains(t, err, "field 'Uplift.Bumps[0].File' contains a path to a file that does not exist 'does-not-exist.txt'")
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
	require.ErrorContains(t, err, "field 'Uplift.Bumps[0].Regex' must be provided when field 'JSON' is missing")
	require.ErrorContains(t, err, "field 'Uplift.Bumps[0].JSON' must be provided when field 'Regex' is missing")
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
	require.ErrorContains(t, err, "field 'Uplift.Bumps[0].Regex[0].Pattern' contains a value that does not meet the minimum expected length of '1'")
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
	require.ErrorContains(t, err, "field 'Uplift.Bumps[0].Regex[0].Count' contains a value that does not meet the minimum expected length of '0'")
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
	require.ErrorContains(t, err, "field 'Uplift.Bumps[0].JSON[0].Path' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateCommitAuthorNameAndEmailEmpty(t *testing.T) {
	cfg := Uplift{
		CommitAuthor: &CommitAuthor{},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.CommitAuthor.Name' must be provided when field 'Email' is missing")
	require.ErrorContains(t, err, "field 'Uplift.CommitAuthor.Email' must be provided when field 'Name' is missing")
}

func TestValidateCommitEmail(t *testing.T) {
	cfg := Uplift{
		CommitAuthor: &CommitAuthor{
			Email: "not-a-valid-email",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.CommitAuthor.Email' contains an invalid email address 'not-a-valid-email'")
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
	require.ErrorContains(t, err, "field 'Uplift.Changelog.Sort' contains a value that is not one of the following [asc desc ASC DESC]")
}

func TestValidateChangeLogExcludeEmpty(t *testing.T) {
	cfg := Uplift{
		Changelog: &Changelog{
			Exclude: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Changelog.Exclude[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateChangelogIncludeEmpty(t *testing.T) {
	cfg := Uplift{
		Changelog: &Changelog{
			Include: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Changelog.Include[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateChangelogAllEmpty(t *testing.T) {
	cfg := Uplift{
		Changelog: &Changelog{},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Changelog.Sort' must be provided when all other fields [Exclude Include] are missing")
	require.ErrorContains(t, err, "field 'Uplift.Changelog.Exclude' must be provided when all other fields [Sort Include] are missing")
	require.ErrorContains(t, err, "field 'Uplift.Changelog.Include' must be provided when all other fields [Sort Exclude] are missing")
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
	require.ErrorContains(t, err, "field 'Uplift.Git.PushOptions[0].Option' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateGitHubNonCompliantURL(t *testing.T) {
	cfg := Uplift{
		GitHub: &GitHub{
			URL: "not-a-url",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.GitHub.URL' contains an invalid url 'not-a-url'")
}

func TestValidateGitLabNonCompliantURL(t *testing.T) {
	cfg := Uplift{
		GitLab: &GitLab{
			URL: "not-a-url",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.GitLab.URL' contains an invalid url 'not-a-url'")
}

func TestValidateGiteaNonCompliantURL(t *testing.T) {
	cfg := Uplift{
		Gitea: &Gitea{
			URL: "not-a-url",
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Gitea.URL' contains an invalid url 'not-a-url'")
}

func TestValidateHooksBeforeEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			Before: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Hooks.Before[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateHooksBeforeBumpEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			BeforeBump: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Hooks.BeforeBump[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateHooksBeforeTagEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			BeforeTag: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Hooks.BeforeTag[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateHooksBeforeChangelogEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			BeforeChangelog: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Hooks.BeforeChangelog[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateHooksAfterEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			After: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Hooks.After[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateHooksAfterBumpEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			AfterBump: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Hooks.AfterBump[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateHooksAfterTagEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			AfterTag: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Hooks.AfterTag[0]' contains a value that does not meet the minimum expected length of '1'")
}

func TestValidateHooksAfterChangelogEmpty(t *testing.T) {
	cfg := Uplift{
		Hooks: &Hooks{
			AfterChangelog: []string{""},
		},
	}

	err := cfg.Validate()
	require.ErrorContains(t, err, "field 'Uplift.Hooks.AfterChangelog[0]' contains a value that does not meet the minimum expected length of '1'")
}
