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

package semver

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBump(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		commit  string
		version string
	}{
		{
			name:    "MajorIncrement",
			tag:     "1.2.3",
			commit:  "refactor!: Lorem ipsum dolor sit amet",
			version: "2.0.0",
		},
		{
			name: "MajorIncrementBreakingChangeFooter",
			tag:  "1.2.3",
			commit: `refactor: Lorem ipsum dolor sit amet

BREAKING CHANGE: Lorem ipsum dolor sit amet`,
			version: "2.0.0",
		},
		{
			name:    "MinorIncrement",
			tag:     "1.2.3",
			commit:  "feat: Lorem ipsum dolor sit amet",
			version: "1.3.0",
		},
		{
			name:    "PatchIncrement",
			tag:     "1.2.3",
			commit:  "fix(db): Lorem ipsum dolor sit amet",
			version: "1.2.4",
		},
		{
			name:    "NoChange",
			tag:     "1.2.3",
			commit:  "chore: Lorem ipsum dolor sit amet",
			version: "1.2.3",
		},
		{
			name:    "KeepsVersionPrefix",
			tag:     "v1.1.1",
			commit:  "fix: Lorem ipsum dolor sit amet",
			version: "v1.1.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MkRepo(t, tt.tag, tt.commit)

			b := NewBumper(io.Discard, BumpOptions{})
			if err := b.Bump(); err != nil {
				t.Errorf("Unexpected error during bump: %s\n", err)
			}

			v := git.LatestTag()

			if v != tt.version {
				t.Errorf("Expected %s but received %s", tt.version, v)
			}
		})
	}
}

func TestBumpDryRun(t *testing.T) {
	v := "1.0.0"
	MkRepo(t, v, "feat: Lorem ipsum dolor sit amet")

	b := NewBumper(io.Discard, BumpOptions{DryRun: true})

	err := b.Bump()
	require.NoError(t, err)

	tag := git.LatestTag()
	assert.Equal(t, v, tag)
}

func TestBumpInvalidVersion(t *testing.T) {
	MkRepo(t, "1.0.B", "feat: Lorem ipsum dolor sit amet")

	b := NewBumper(io.Discard, BumpOptions{})
	err := b.Bump()

	require.Error(t, err)
}

func TestBumpFirstVersion(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: Lorem ipsum dolor sit amet")

	b := NewBumper(io.Discard, BumpOptions{Config: config.Uplift{FirstVersion: "0.1.0"}})
	err := b.Bump()
	require.NoError(t, err)

	v := git.LatestTag()
	assert.Equal(t, "0.1.0", v)
}

func TestBumpEmptyRepo(t *testing.T) {
	git.InitRepo(t)

	b := NewBumper(io.Discard, BumpOptions{})
	err := b.Bump()

	require.NoError(t, err)
}

func TestBumpNotGitRepo(t *testing.T) {
	git.MkTmpDir(t)

	b := NewBumper(io.Discard, BumpOptions{})
	err := b.Bump()

	require.Error(t, err)
	assert.Error(t, err, "current directory must be a git repo")
}

func TestBumpAlwaysUseLatestCommit(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommits(t,
		"feat: Lorem ipsum dolor sit amet",
		"fix: Lorem ipsum dolor sit amet",
		"docs: Lorem ipsum dolor sit amet")

	b := NewBumper(io.Discard, BumpOptions{})
	err := b.Bump()

	require.NoError(t, err)
	assert.Equal(t, "", git.LatestTag())
}

func TestBumpWithAnnotatedTag(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommitAndTag(t, "0.1.0", "feat: Lorem ipsum dolor sit amet")

	b := NewBumper(io.Discard, BumpOptions{
		Config: config.Uplift{
			AnnotatedTags: true,
			CommitMessage: "this is an annotated tag",
		},
	})
	err := b.Bump()

	require.NoError(t, err)
	assert.Equal(t, "0.2.0", git.LatestTag())

	out, _ := git.Clean(git.Run("for-each-ref", "refs/tags/0.2.0",
		"--format='%(contents)'"))

	assert.Contains(t, out, "this is an annotated tag")
}

func TestBumpFile(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		commit   string
		regex    string
		count    int
		content  string
		expected string
	}{
		{
			name:     "SingleLineAtEnd",
			tag:      "v1.0.0",
			commit:   "fix: Lorem ipsum dolor sit amet",
			regex:    `version $VERSION`,
			content:  "file is currently at version v1.0.0",
			expected: "file is currently at version v1.0.1",
		},
		{
			name:     "SingleLineAtStart",
			tag:      "1.0.0",
			commit:   "feat: Lorem ipsum dolor sit amet",
			regex:    "version $VERSION",
			content:  "version 1.0.0 contains many new changes",
			expected: "version 1.1.0 contains many new changes",
		},
		{
			name: "SingleLineWrapped",
			tag:  "1.0.0",
			commit: `feat: Lorem ipsum dolor sit amet
BREAKING CHANGE: Lorem ipsum dolor sit amet`,
			regex:    "version: $VERSION",
			content:  "file is at [version: 1.0.0], updated on 24th April 2021",
			expected: "file is at [version: 2.0.0], updated on 24th April 2021",
		},
		{
			name:   "Multiline",
			tag:    "2.4.0",
			commit: "feat: Lorem ipsum dolor sit amet",
			regex:  "version: $VERSION",
			content: `this is a test to see if the version is updated in
		a file with multiple lines. This is now at version: 2.4.0.
		Updated on 10th May 2021`,
			expected: `this is a test to see if the version is updated in
		a file with multiple lines. This is now at version: 2.5.0.
		Updated on 10th May 2021`,
		},
		{
			name:   "MultilineReplaceAll",
			tag:    "1.0.0",
			commit: "fix!: Lorem ipsum dolor sit amet",
			regex:  "version $VERSION",
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
			name:   "MultilineReplaceSingle",
			tag:    "1.0.0",
			commit: "fix: Lorem ipsum dolor sit amet",
			regex:  "version: $VERSION",
			count:  1,
			content: `version: 1.0.0
		version: 1.0.0
		version: 1.0.0`,
			expected: `version: 1.0.1
		version: 1.0.0
		version: 1.0.0`,
		},
		{
			name:   "MultilineReplaceMultiple",
			tag:    "v1.0.0",
			commit: "feat!: Lorem ipsum dolor sit amet",
			regex:  "version: $VERSION",
			count:  2,
			content: `version: v1.0.0
		version: v1.0.0
		version: v1.0.0`,
			expected: `version: v2.0.0
		version: v2.0.0
		version: v1.0.0`,
		},
		{
			name:   "RegexSpecialCharacters",
			tag:    "1.0.0",
			commit: "feat!: Lorem ipsum dolor sit amet",
			regex:  `\s{2}version: $VERSION`,
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
			MkRepo(t, tt.tag, tt.commit)

			// Generate file to bump
			path := WriteFile(t, tt.content)

			opts := BumpOptions{
				Config: config.Uplift{
					Bumps: []config.Bump{
						{
							File:  path,
							Regex: tt.regex,
							Count: tt.count,
						},
					},
				},
			}

			b := NewBumper(io.Discard, opts)
			if err := b.Bump(); err != nil {
				t.Errorf("Unexpected error during bump: %s\n", err)
			}

			actual := ReadFile(t, path)

			if actual != tt.expected {
				t.Errorf("Expected:\n%s\nbut received:\n%s", tt.expected, actual)
			}
		})
	}
}

func TestBumpFileSemanticVersionOnly(t *testing.T) {
	MkRepo(t, "v0.1.0", "feat: Lorem ipsum dolor sit amet")

	file := "version: 0.1.0"
	path := WriteFile(t, file)

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:   path,
					Regex:  "version: $VERSION",
					SemVer: true,
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	actual := ReadFile(t, path)
	assert.Equal(t, "version: 0.2.0", actual)
}

func TestBumpFileDryRun(t *testing.T) {
	MkRepo(t, "0.1.0", "fix: Lorem ipsum dolor sit amet")

	file := "version: 0.1.0"
	path := WriteFile(t, file)

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  path,
					Regex: "version: $VERSION",
				},
			},
		},
		DryRun: true,
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	actual := ReadFile(t, path)
	assert.Equal(t, file, actual)
}

func TestBumpFileDefaultCommitMessage(t *testing.T) {
	MkRepo(t, "0.1.0", "fix: Lorem ipsum dolor sit amet")

	file := "version: 0.1.0"
	path := WriteFile(t, file)

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  path,
					Regex: "version: $VERSION",
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	commit, err := git.LatestCommit()
	require.NoError(t, err)

	assert.Equal(t, "uplift", commit.Author)
	assert.Equal(t, "uplift@test.com", commit.Email)
	assert.Equal(t, "ci(bump): bumped version to 0.1.1", commit.Message)
}

func TestBumpFileWithCustomisedCommit(t *testing.T) {
	MkRepo(t, "0.1.0", "fix: Lorem ipsum dolor sit amet")

	file := "version: 0.1.0"
	path := WriteFile(t, file)

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  path,
					Regex: "version: $VERSION",
				},
			},
			CommitMessage: "chore: Lorem ipsum dolor sit amet",
			CommitAuthor: config.CommitAuthor{
				Name:  "joe.bloggs",
				Email: "joe.bloggs@gmail.com",
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	commit, err := git.LatestCommit()
	require.NoError(t, err)

	assert.Equal(t, "joe.bloggs", commit.Author)
	assert.Equal(t, "joe.bloggs@gmail.com", commit.Email)
	assert.Equal(t, "chore: Lorem ipsum dolor sit amet", commit.Message)
}

func TestBumpFileFirstTagMatchesVersionInFile(t *testing.T) {
	git.InitRepo(t)
	git.EmptyCommit(t, "feat: Lorem ipsum dolor sit amet")

	file := "version: 0.1.0"
	path := WriteFile(t, file)

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  path,
					Regex: "version: $VERSION",
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	commit, err := git.LatestCommit()
	require.NoError(t, err)

	assert.Equal(t, "feat: Lorem ipsum dolor sit amet", commit.Message)
}

func TestBumpFileOnMissingFile(t *testing.T) {
	MkRepo(t, "0.1.0", "fix: Lorem ipsum dolor sit amet")

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  "file.txt",
					Regex: "version: $VERSION",
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.Error(t, err, "open file.txt: no such file or directory")
}

func TestBumpMultipleFiles(t *testing.T) {
	MkRepo(t, "0.1.0", "fix: Lorem ipsum dolor sit amet")

	contents := "version: 0.1.0"
	file1 := WriteFile(t, contents)
	file2 := WriteFile(t, contents)

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  file1,
					Regex: "version: $VERSION",
				},
				{
					File:  file2,
					Regex: "version: $VERSION",
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	expected := "version: 0.1.1"
	actual1 := ReadFile(t, file1)
	assert.Equal(t, expected, actual1)

	actual2 := ReadFile(t, file2)
	assert.Equal(t, expected, actual2)
}

func TestBumpFileNonMatchingRegex(t *testing.T) {
	MkRepo(t, "0.1.0", "fix: Lorem ipsum dolor sit amet")

	path := WriteFile(t, "version: 0.1.0")

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  path,
					Regex: "noMatch: $VERSION",
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.Error(t, err, "no version matched in file")
}

func MkRepo(t *testing.T, tag, commit string) {
	t.Helper()
	git.InitRepo(t)

	err := git.Tag(tag)
	require.NoError(t, err)

	git.EmptyCommit(t, commit)
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

func TestBumpMavenPom(t *testing.T) {
	MkRepo(t, "0.1.0", "feat: Lorem ipsum dolor sit amet")

	path := WriteFile(t, `<?xml version="1.0" encoding="UTF-8"?>
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

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  path,
					Regex: `\s{4}<version>$VERSION</version>`,
					Count: 1,
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	actual := ReadFile(t, path)
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

func TestBumpHelmChart(t *testing.T) {
	MkRepo(t, "0.1.0", "fix: Lorem ipsum dolor sit amet")

	path := WriteFile(t, `apiVersion: v2
name: test-chart
description: This is a test chart
version: 0.1.0
appVersion: 0.1.0`)

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  path,
					Regex: "version: $VERSION",
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	actual := ReadFile(t, path)
	assert.Equal(t, `apiVersion: v2
name: test-chart
description: This is a test chart
version: 0.1.1
appVersion: 0.1.0`, actual)
}

func TestBumpPackageJson(t *testing.T) {
	MkRepo(t, "0.1.0", "feat!: Lorem ipsum dolor sit amet")

	path := WriteFile(t, `{
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

	opts := BumpOptions{
		Config: config.Uplift{
			Bumps: []config.Bump{
				{
					File:  path,
					Regex: `"version": "$VERSION"`,
				},
			},
		},
	}

	b := NewBumper(io.Discard, opts)
	err := b.Bump()
	require.NoError(t, err)

	actual := ReadFile(t, path)
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
