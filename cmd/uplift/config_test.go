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

package main

import (
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			name:     "DotUpliftYml",
			filename: ".uplift.yml",
		},
		{
			name:     "DotUpliftYaml",
			filename: ".uplift.yaml",
		},
		{
			name:     "UpliftYml",
			filename: "uplift.yml",
		},
		{
			name:     "UpliftYaml",
			filename: "uplift.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gittest.InitRepository(t)
			gittest.TempFile(t, tt.filename, "annotatedTags: true")

			cfg, err := loadConfig(currentWorkingDir)

			require.NoError(t, err)
			require.True(t, cfg.AnnotatedTags)
		})
	}
}

func TestLoadConfig_Malformed(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, ".uplift.yml", "firstV")

	_, err := loadConfig(currentWorkingDir)
	assert.Error(t, err)
}

func TestLoadConfig_NotExists(t *testing.T) {
	gittest.InitRepository(t)

	_, err := loadConfig(currentWorkingDir)
	assert.NoError(t, err)
}

func TestLoadConfig_CustomLocation(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, "custom/.uplift.yml", "annotatedTags: true")

	cfg, err := loadConfig("custom")
	assert.NoError(t, err)
	require.True(t, cfg.AnnotatedTags)
}
