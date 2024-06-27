package main

import (
	"os"
	"path/filepath"

	"github.com/gembaadvantage/uplift/internal/config"
)

var files = [4]string{".uplift.yml", ".uplift.yaml", "uplift.yml", "uplift.yaml"}

const (
	currentWorkingDir = "."
)

func loadConfig(dir string) (config.Uplift, error) {
	for _, file := range files {
		cfg, err := config.Load(filepath.Join(dir, file))

		// If the file doesn't exist, try another, until the array is exhausted
		if err != nil && os.IsNotExist(err) {
			continue
		}

		return cfg, err
	}

	return config.Uplift{}, nil
}
