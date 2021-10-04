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
	"strconv"
)

const (
	// Pattern defines the regular expression for matching a semantic version
	// as a raw sequence of characters
	Pattern = `(v?)(\d+)\.(\d+)\.(\d+)`
)

var (
	// Regex for parsing a matching a semantic version
	Regex = regexp.MustCompile(Pattern)
)

// Version provides a less strict implementation of a semantic version
// and supports the optional use of a prefix, which is a common standard with
// Git tags
type Version struct {
	Prefix string
	Patch  uint64
	Minor  uint64
	Major  uint64
	Raw    string
}

// Parse ...
func Parse(ver string) (Version, error) {
	if m := Regex.FindStringSubmatch(ver); len(m) > 4 {
		return Version{
			Prefix: m[1],
			Major:  toUint64(m[2]),
			Minor:  toUint64(m[3]),
			Patch:  toUint64(m[4]),
			Raw:    ver,
		}, nil
	}

	return Version{}, errors.New("unrecognised semantic version format")
}

func toUint64(d string) uint64 {
	v, _ := strconv.Atoi(d)
	return uint64(v)
}

// String outputs the unparsed semantic version
func (v Version) String() string {
	return v.Raw
}
