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
	"errors"
	"regexp"

	semv "github.com/Masterminds/semver"
)

const (
	// Pattern defines the regular expression for matching a semantic version.
	// Taken directly from github.com/Masterminds/semver
	Pattern = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
		`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
		`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

	// Token defines a constant that can be used to perform a string replacement
	// in a consistent manner. Will be replaced with template support in future
	Token = "$VERSION"
)

var (
	// Regex for matching a semantic version
	Regex = regexp.MustCompile(Pattern)
)

// Version provides a less strict implementation of a semantic version
// by supporting an optional use of a 'v' prefix
type Version struct {
	Prefix     string
	Patch      int64
	Minor      int64
	Major      int64
	Prerelease string
	Metadata   string
	Raw        string
}

// Parse a semantic version
func Parse(ver string) (Version, error) {
	v, err := semv.NewVersion(ver)
	if err != nil {
		return Version{}, err
	}

	// Detect and capture optionally supported prefix
	prefix := ""
	if ver[0] == 'v' {
		prefix = "v"
	}

	return Version{
		Prefix:     prefix,
		Major:      v.Major(),
		Minor:      v.Minor(),
		Patch:      v.Patch(),
		Prerelease: v.Prerelease(),
		Metadata:   v.Metadata(),
		Raw:        ver,
	}, nil
}

// ParsePrerelease attempts to parse a prerelease suffix. Supports
// a prerelease suffix with and without a leading '-'
func ParsePrerelease(pre string) (string, string, error) {
	if pre == "" {
		return "", "", errors.New("prerelease suffix is blank")
	}

	// Has prefix been provided
	i := 0
	if pre[0] == '-' {
		i = 1
	}

	v, err := Parse("1.0.0-" + pre[i:])
	if err != nil {
		return "", "", errors.New("invalid semantic prerelease suffix")
	}

	return v.Prerelease, v.Metadata, nil
}

// String outputs the unparsed semantic version
func (v Version) String() string {
	return v.Raw
}
