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

package semver

import "regexp"

// Increment defines the different types of increment that can be performed
// against a semantic version
type Increment string

const (
	// NoIncrement represents no increment change to a semantic version
	NoIncrement Increment = "None"
	// PatchIncrement represents a patch increment (1.0.x) to a semantic version
	PatchIncrement Increment = "Patch"
	// MinorIncrement represents a minor increment (1.x.0) to a semantic version
	MinorIncrement Increment = "Minor"
	// MajorIncrement represents a major increment (x.0.0) to a semantic version
	MajorIncrement Increment = "Major"
)

var (
	breakingBang = regexp.MustCompile(`(?im).*(\w+)(\(.*\))?!:.*`)
	breaking     = regexp.MustCompile("(?im).*BREAKING CHANGE:.*")
	feature      = regexp.MustCompile(`(?im).*feat(\(.*\))?:.*`)
	fix          = regexp.MustCompile(`(?im).*fix(\(.*\))?:.*`)
)

// TODO: return line that triggered the increment

// ParseLog will identify the maximum semantic increment by parsing the commit
// log against the conventional commit standards defined, @see:
// https://www.conventionalcommits.org/en/v1.0.0/
func ParseLog(log string) Increment {
	if breakingBang.MatchString(log) || breaking.MatchString(log) {
		return MajorIncrement
	}

	if feature.MatchString(log) {
		return MinorIncrement
	}

	if fix.MatchString(log) {
		return PatchIncrement
	}

	return NoIncrement
}
