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

// TODO: move files to regex_test.go

package bump

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	commit = git.CommitDetails{
		Author:  "joe.bloggs",
		Email:   "joe.bloggs@example.com",
		Message: "dummy commit",
	}
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		nextVer  string
		regex    string
		count    int
		content  string
		expected string
	}{
		{
			name:     "SingleLineAtEnd",
			nextVer:  "v1.0.1",
			regex:    "version $VERSION",
			content:  "file is currently at version v1.0.0",
			expected: "file is currently at version v1.0.1",
		},
		{
			name:     "SingleLineAtStart",
			nextVer:  "1.1.0",
			regex:    "$VERSION",
			content:  "1.0.0 contains many new changes",
			expected: "1.1.0 contains many new changes",
		},
		{
			name:     "SingleWithPrerelease",
			nextVer:  "0.3.1-beta.1+ade3f12",
			regex:    "version $VERSION",
			content:  "file is at version 0.3.0",
			expected: "file is at version 0.3.1-beta.1+ade3f12",
		},
		{
			name:     "SingleLineWrapped",
			nextVer:  "2.0.0",
			regex:    "version: $VERSION",
			content:  "file is at [version: 1.0.0], updated on 24th April 2021",
			expected: "file is at [version: 2.0.0], updated on 24th April 2021",
		},
		{
			name:    "Multiline",
			nextVer: "2.5.0",
			regex:   "version: $VERSION",
			content: `this is a test to see if the version is updated in
		a file with multiple lines. This is now at version: 2.4.0.
		Updated on 10th May 2021`,
			expected: `this is a test to see if the version is updated in
		a file with multiple lines. This is now at version: 2.5.0.
		Updated on 10th May 2021`,
		},
		{
			name:    "MultilineReplaceAll",
			nextVer: "2.0.0",
			regex:   "version $VERSION",
			content: `this is a test to check all occurrences of the version are replaced.
		From here: version 1.0.0
		And here: version 1.0.0
		And also here: version 1.0.0 version 1.0.0
		version 1.0.0`,
			expected: `this is a test to check all occurrences of the version are replaced.
		From here: version 2.0.0
		And here: version 2.0.0
		And also here: version 2.0.0 version 2.0.0
		version 2.0.0`,
		},
		{
			name:    "MultilineReplaceSingle",
			nextVer: "1.0.1",
			regex:   "version: $VERSION",
			count:   1,
			content: `version: 1.0.0
		version: 1.0.0
		version: 1.0.0`,
			expected: `version: 1.0.1
		version: 1.0.0
		version: 1.0.0`,
		},
		{
			name:    "MultilineReplaceMultiple",
			nextVer: "v2.0.0",
			regex:   "version: $VERSION",
			count:   2,
			content: `version: v1.0.0
		version: v1.0.0
		version: v1.0.0`,
			expected: `version: v2.0.0
		version: v2.0.0
		version: v1.0.0`,
		},
		{
			name:    "RegexSpecialCharacters",
			nextVer: "2.0.0",
			regex:   `\s{2}version: $VERSION`,
			content: `this is a test:
		  version: 1.0.0
		version: 1.0.0
		  version: 1.0.0
		version: 1.0.0`,
			expected: `this is a test:
		  version: 2.0.0
		version: 1.0.0
		  version: 2.0.0
		version: 1.0.0`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			git.InitRepo(t)
			path := WriteFile(t, tt.content)

			ctx := &context.Context{
				NextVersion: semver.Version{
					Raw: tt.nextVer,
				},
				CommitDetails: commit,
				Config: config.Uplift{
					Bumps: []config.Bump{
						{
							File: path,
							Regex: []config.RegexBump{
								{
									Pattern: tt.regex,
									Count:   tt.count,
								},
							},
						},
					},
				},
			}

			err := Task{}.Run(ctx)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			actual := ReadFile(t, path)

			if actual != tt.expected {
				t.Errorf("Expected:\n%s\nbut received:\n%s", tt.expected, actual)
			}
		})
	}
}

func TestRun_ForceSemanticVersion(t *testing.T) {
	git.InitRepo(t)
	path := WriteFile(t, "version: 0.1.0")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "v0.2.0",
		},
		CommitDetails: commit,
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: path,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
							SemVer:  true,
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, path)
	assert.Equal(t, "version: 0.2.0", actual)
}

func TestRun_DryRun(t *testing.T) {
	path := WriteFile(t, "version: 0.1.0")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: path,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
						},
					},
				},
			},
		},
		DryRun: true,
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, path)
	assert.Equal(t, "version: 0.1.0", actual)
}

func TestRun_FileDoesNotExist(t *testing.T) {
	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.2.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: "missing.txt",
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	assert.Error(t, err)
}

func TestRun_MultipleFiles(t *testing.T) {
	git.InitRepo(t)

	contents := "version: 0.1.0"
	file1 := WriteFile(t, contents)
	file2 := WriteFile(t, contents)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		CommitDetails: commit,
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file1,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
						},
					},
				},
				{
					File: file2,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	expected := fmt.Sprintf("version: %s", ctx.NextVersion.Raw)
	actual1 := ReadFile(t, file1)
	assert.Equal(t, expected, actual1)

	actual2 := ReadFile(t, file2)
	assert.Equal(t, expected, actual2)
}

func TestRun_NonMatchingRegex(t *testing.T) {
	git.InitRepo(t)
	file := WriteFile(t, "version: 0.1.0")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: "noMatch: $VERSION",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	assert.EqualError(t, err, "no version matched in file")
}

func TestRun_NoBumpConfig(t *testing.T) {
	err := Task{}.Run(&context.Context{})
	assert.NoError(t, err)
}

func TestRun_NotGitRepository(t *testing.T) {
	git.MkTmpDir(t)
	file := WriteFile(t, "version: 0.1.0")

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	assert.EqualError(t, err, "fatal: not a git repository (or any of the parent directories): .git")
}

func TestRun_NextVersionMatchesExistingVersion(t *testing.T) {
	git.InitRepo(t)
	file := WriteFile(t, "version: 0.1.0")

	efi, _ := os.Stat(file)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.0",
		},
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	// Check that the file has not been modified
	afi, _ := os.Stat(file)
	assert.Equal(t, efi.ModTime(), afi.ModTime())
}

func TestRun_MalformedRegexError(t *testing.T) {
	git.MkTmpDir(t)
	file := WriteFile(t, "version: 0.1.0")

	ctx := &context.Context{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: "[",
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	assert.Error(t, err)
}

func WriteFile(t *testing.T, s string) string {
	t.Helper()

	current, err := os.Getwd()
	require.NoError(t, err)

	file, err := ioutil.TempFile(current, "*")
	require.NoError(t, err)

	_, err = file.WriteString(s)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	t.Cleanup(func() {
		require.NoError(t, os.Remove(file.Name()))
	})

	return file.Name()
}

func ReadFile(t *testing.T, path string) string {
	t.Helper()

	b, err := ioutil.ReadFile(path)
	require.NoError(t, err)

	return string(b)
}

func TestRun_MavenPom(t *testing.T) {
	git.InitRepo(t)

	file := WriteFile(t, `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>2.3.5.RELEASE</version>
        <relativePath/>
    </parent>
    <groupId>com.test.organisation</groupId>
    <artifactId>test-service</artifactId>
    <version>0.1.0</version>
    <name>test-service</name>
    <description>This is a test service</description>

	<dependencies>
        <dependency>
            <groupId>org.apache.commons</groupId>
            <artifactId>commons-lang3</artifactId>
            <version>3.12.0</version>
        </dependency>
    </dependencies>
</project>`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.2.0",
		},
		CommitDetails: commit,
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: `\s{4}<version>$VERSION</version>`,
							Count:   1,
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, file)
	assert.Equal(t, `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>2.3.5.RELEASE</version>
        <relativePath/>
    </parent>
    <groupId>com.test.organisation</groupId>
    <artifactId>test-service</artifactId>
    <version>0.2.0</version>
    <name>test-service</name>
    <description>This is a test service</description>

	<dependencies>
        <dependency>
            <groupId>org.apache.commons</groupId>
            <artifactId>commons-lang3</artifactId>
            <version>3.12.0</version>
        </dependency>
    </dependencies>
</project>`, actual)
}

func TestRun_HelmChart(t *testing.T) {
	git.InitRepo(t)

	file := WriteFile(t, `apiVersion: v2
name: test-chart
description: This is a test chart
version: 0.1.0
appVersion: 0.1.0`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "0.1.1",
		},
		CommitDetails: commit,
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: "version: $VERSION",
							Count:   1,
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, file)
	assert.Equal(t, `apiVersion: v2
name: test-chart
description: This is a test chart
version: 0.1.1
appVersion: 0.1.0`, actual)
}

func TestRun_PackageJson(t *testing.T) {
	git.InitRepo(t)

	file := WriteFile(t, `{
  "name": "test",
  "version": "0.1.0",
  "bin": {
    "test": "bin/test.js"
  },
  "scripts": {
    "build": "tsc",
  },
  "devDependencies": {
    "typescript": "~3.7.2"
  },
  "dependencies": {}
}`)

	ctx := &context.Context{
		NextVersion: semver.Version{
			Raw: "1.0.0",
		},
		CommitDetails: commit,
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File: file,
					Regex: []config.RegexBump{
						{
							Pattern: `"version": "$VERSION"`,
							Count:   1,
						},
					},
				},
			},
		},
	}

	err := Task{}.Run(ctx)
	require.NoError(t, err)

	actual := ReadFile(t, file)
	assert.Equal(t, `{
  "name": "test",
  "version": "1.0.0",
  "bin": {
    "test": "bin/test.js"
  },
  "scripts": {
    "build": "tsc",
  },
  "devDependencies": {
    "typescript": "~3.7.2"
  },
  "dependencies": {}
}`, actual)
}
