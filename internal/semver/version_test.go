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
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_SemVer(t *testing.T) {
	v, err := Parse("1.22.1-beta.1+12345")
	require.NoError(t, err)

	assert.Equal(t, "", v.Prefix)
	assert.Equal(t, int64(1), v.Major)
	assert.Equal(t, int64(22), v.Minor)
	assert.Equal(t, int64(1), v.Patch)
	assert.Equal(t, "beta.1", v.Prerelease)
	assert.Equal(t, "12345", v.Metadata)
}

func TestParse_WithPrefix(t *testing.T) {
	v, err := Parse("v0.3.1")
	require.NoError(t, err)

	assert.Equal(t, "v", v.Prefix)
	assert.Equal(t, int64(0), v.Major)
	assert.Equal(t, int64(3), v.Minor)
	assert.Equal(t, int64(1), v.Patch)
	assert.Equal(t, "", v.Prerelease)
	assert.Equal(t, "", v.Metadata)
}

func TestParse_InvalidSemVer(t *testing.T) {
	_, err := Parse("V1.0.0")
	require.Error(t, err)
}

func TestString_ReturnsRaw(t *testing.T) {
	v := Version{Raw: "1.0.0-beta.1"}

	var buf bytes.Buffer
	fmt.Fprint(&buf, v.String())

	assert.Equal(t, "1.0.0-beta.1", buf.String())
}

func TestParsePrerelease(t *testing.T) {
	tests := []struct {
		name       string
		prerelease string
		pre        string
		meta       string
	}{
		{
			name:       "WithLeadingHyphen",
			prerelease: "-beta.1+a2sd3ef",
			pre:        "beta.1",
			meta:       "a2sd3ef",
		},
		{
			name:       "NoLeadingHyphen",
			prerelease: "beta.1+a2sd3ef",
			pre:        "beta.1",
			meta:       "a2sd3ef",
		},
		{
			name:       "NoMetadata",
			prerelease: "beta.1",
			pre:        "beta.1",
			meta:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre, meta, err := ParsePrerelease(tt.prerelease)

			require.NoError(t, err)
			require.Equal(t, tt.pre, pre)
			require.Equal(t, tt.meta, meta)
		})
	}
}

func TestParsePrerelease_Empty(t *testing.T) {
	_, _, err := ParsePrerelease("")

	assert.EqualError(t, err, "prerelease suffix is blank")
}

func TestParsePrerelease_Invalid(t *testing.T) {
	_, _, err := ParsePrerelease("-#")

	assert.EqualError(t, err, "invalid semantic prerelease suffix")
}
