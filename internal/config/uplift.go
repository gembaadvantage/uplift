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
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// Uplift defines the root configuration of the application
type Uplift struct {
	Git           Git          `yaml:"git"`
	FirstVersion  string       `yaml:"firstVersion"`
	Bumps         []Bump       `yaml:"bumps"`
	CommitMessage string       `yaml:"commitMessage"`
	CommitAuthor  CommitAuthor `yaml:"commitAuthor"`
	AnnotatedTags bool         `yaml:"annotatedTags"`
	Changelog     Changelog    `yaml:"changelog"`
}

// Git defines configuration for how uplift interacts with git
type Git struct {
	IgnoreDetached bool `yaml:"ignoreDetached"`
	IgnoreShallow  bool `yaml:"ignoredShallow"`
}

// Bump defines configuration for bumping individual files based
// on the new calculated semantic version number
type Bump struct {
	File  string      `yaml:"file"`
	Regex []RegexBump `yaml:"regex"`
	JSON  []JSONBump  `yaml:"json"`
}

// RegexBump defines configuration for bumping a file based on
// a given regex pattern
type RegexBump struct {
	Pattern string `yaml:"pattern"`
	Count   int    `yaml:"count"`
	SemVer  bool   `yaml:"semver"`
}

// JSONBump defines configuration for bumping a file based on a
// given JSON path. Path syntax is based on the github.com/tidwall/sjson
// library
type JSONBump struct {
	Path   string `yaml:"path"`
	SemVer bool   `yaml:"semver"`
}

// CommitAuthor defines configuration about the author of a git commit
type CommitAuthor struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

// Changelog defines configuration for generating a changelog of the latest
// semantic version based release
type Changelog struct {
	Sort    string   `yaml:"sort"`
	Exclude []string `yaml:"exclude"`
}

// Load the YAML config file
func Load(f string) (Uplift, error) {
	fh, err := os.Open(f)
	if err != nil {
		return Uplift{}, err
	}
	defer fh.Close()

	// Read the contents of the file in one go
	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return Uplift{}, err
	}

	var cfg Uplift
	err = yaml.UnmarshalStrict(data, &cfg)
	return cfg, err
}
