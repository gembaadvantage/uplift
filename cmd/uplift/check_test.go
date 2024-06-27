package main

import (
	"os"
	"testing"

	"github.com/purpleclay/gitz/gittest"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, ".uplift.yml", `commitAuthor:
  name: joe.bloggs
  email: joe.bloggs@example.com
`)

	checkCmd := newCheckCmd(&globalOptions{}, os.Stdout)
	err := checkCmd.Execute()

	assert.NoError(t, err)
}

func TestCheck_InvalidConfig(t *testing.T) {
	gittest.InitRepository(t)
	gittest.TempFile(t, ".uplift.yml", `bumps:
  - file: text.txt
    regex:
      - pattern: ""
`)

	checkCmd := newCheckCmd(&globalOptions{}, os.Stdout)
	err := checkCmd.Execute()

	assert.Error(t, err)
}
