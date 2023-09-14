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
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Uplift defines the root configuration of the application
type Uplift struct {
	AnnotatedTags bool          `yaml:"annotatedTags"`
	Bumps         []Bump        `yaml:"bumps" validate:"omitempty,dive"`
	CommitAuthor  *CommitAuthor `yaml:"commitAuthor" validate:"omitempty"`
	CommitMessage string        `yaml:"commitMessage"`
	Changelog     *Changelog    `yaml:"changelog" validate:"omitempty"`
	Git           *Git          `yaml:"git" validate:"omitempty"`
	Gitea         *Gitea        `yaml:"gitea" validate:"omitempty"`
	GitHub        *GitHub       `yaml:"github" validate:"omitempty"`
	GitLab        *GitLab       `yaml:"gitlab" validate:"omitempty"`
	Hooks         *Hooks        `yaml:"hooks" validate:"omitempty"`
	Env           []string      `yaml:"env" validate:"dive,min=1"`
}

// Bump defines configuration for bumping individual files based
// on the new calculated semantic version number
type Bump struct {
	File  string      `yaml:"file" validate:"min=1,file"`
	Regex []RegexBump `yaml:"regex" validate:"required_without=JSON,dive"`
	JSON  []JSONBump  `yaml:"json" validate:"required_without=Regex,dive"`
}

// RegexBump defines configuration for bumping a file based on
// a given regex pattern
type RegexBump struct {
	Pattern string `yaml:"pattern" validate:"min=1"`
	Count   int    `yaml:"count" validate:"min=0"`
	SemVer  bool   `yaml:"semver"`
}

// JSONBump defines configuration for bumping a file based on a
// given JSON path. Path syntax is based on the github.com/tidwall/sjson
// library
type JSONBump struct {
	Path   string `yaml:"path" validate:"min=1"`
	SemVer bool   `yaml:"semver"`
}

// CommitAuthor defines configuration about the author of a git commit
type CommitAuthor struct {
	Name  string `yaml:"name" validate:"required_without=Email,min=1"`
	Email string `yaml:"email" validate:"required_without=Name,email"`
}

// Changelog defines configuration for generating a changelog of the latest
// semantic version based release
type Changelog struct {
	Sort           string   `yaml:"sort" validate:"required_without_all=Exclude Include,oneof=asc desc ASC DESC"`
	Exclude        []string `yaml:"exclude" validate:"required_without_all=Sort Include,dive,min=1"`
	Include        []string `yaml:"include" validate:"required_without_all=Sort Exclude,dive,min=1"`
	Multiline      bool     `yaml:"multiline"`
	SkipPrerelease bool     `yaml:"skipPrerelease"`
}

// Git defines configuration for how uplift interacts with git
type Git struct {
	IgnoreDetached   bool            `yaml:"ignoreDetached"`
	IgnoreShallow    bool            `yaml:"ignoreShallow"`
	PushOptions      []GitPushOption `yaml:"pushOptions" validate:"dive"`
	IncludeArtifacts []string        `yaml:"includeArtifacts" validate:"omitempty,dive,min=1,file"`
}

// GitPushOption provides a way of supplying additional options to
// git push commands
type GitPushOption struct {
	Option     string `yaml:"option" validate:"min=1"`
	SkipBranch bool   `yaml:"skipBranch"`
	SkipTag    bool   `yaml:"skipTag"`
}

type gitPushOption GitPushOption

// UnmarshalYAML defines a custom YAML unmarshal for a [config.GitPushOption]
func (o *GitPushOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err == nil {
		o.Option = str
		return nil
	}

	var opt gitPushOption
	if err := unmarshal(&opt); err != nil {
		return err
	}

	o.Option = opt.Option
	o.SkipBranch = opt.SkipBranch
	o.SkipTag = opt.SkipTag

	return nil
}

// Gitea defines custom configuration for accessing a self-hosted Gitea instance
type Gitea struct {
	URL string `yaml:"url" validate:"url"`
}

// GitHub defines custom configuration for accessing a GitHub enterprise instance
type GitHub struct {
	URL string `yaml:"url" validate:"url"`
}

// GitLab defines custom configuration for accessing a GitLab enterprise or
// self-hosted instance
type GitLab struct {
	URL string `yaml:"url" validate:"url"`
}

// Hooks define custom configuration for entry points before any uplift
// workflow. These entry points can be used to execute any custom shell
// commands or scripts
type Hooks struct {
	Before          []string `yaml:"before" validate:"dive,min=1"`
	BeforeBump      []string `yaml:"beforeBump" validate:"dive,min=1"`
	BeforeTag       []string `yaml:"beforeTag" validate:"dive,min=1"`
	BeforeChangelog []string `yaml:"beforeChangelog" validate:"dive,min=1"`
	After           []string `yaml:"after" validate:"dive,min=1"`
	AfterBump       []string `yaml:"afterBump" validate:"dive,min=1"`
	AfterTag        []string `yaml:"afterTag" validate:"dive,min=1"`
	AfterChangelog  []string `yaml:"afterChangelog" validate:"dive,min=1"`
}

// Load the YAML config file
func Load(f string) (Uplift, error) {
	fh, err := os.Open(f)
	if err != nil {
		return Uplift{}, err
	}
	defer fh.Close()

	// Read the contents of the file in one go
	data, err := io.ReadAll(fh)
	if err != nil {
		return Uplift{}, err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var cfg Uplift
	err = decoder.Decode(&cfg)

	return cfg, err
}

// Validate the existing config, ensuring all values meet expected
// criteria. This also ensures any config file aligns with the
// schema [https://upliftci.dev/static/schema.json]
func (c Uplift) Validate() error {
	if err := validator.New().Struct(c); err != nil {
		var errMsg strings.Builder
		errMsg.WriteString("uplift configuration contains validation errors. Please fix before proceeding:\n\n")

		for _, err := range err.(validator.ValidationErrors) {
			var reason string
			switch err.ActualTag() {
			case "min":
				reason = fmt.Sprintf("contains a value that does not meet the minimum expected length of '%s'\n", err.Param())
			case "file":
				reason = fmt.Sprintf("contains a path to a file that does not exist '%v'\n", err.Value())
			case "url":
				reason = fmt.Sprintf("contains an invalid url '%s'\n", err.Value())
			case "email":
				reason = fmt.Sprintf("contains an invalid email address '%s'\n", err.Value())
			case "oneof":
				reason = fmt.Sprintf("contains a value that is not one of the following [%s]\n", err.Param())
			case "required_without":
				reason = fmt.Sprintf("must be provided when field '%s' is missing\n", err.Param())
			case "required_without_all":
				reason = fmt.Sprintf("must be provided when all other fields [%s] are missing\n", err.Param())
			}

			errMsg.WriteString(fmt.Sprintf(" field '%s' ", err.Namespace()))
			errMsg.WriteString(reason)
		}

		return errors.New(errMsg.String())
	}

	return nil
}
