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
