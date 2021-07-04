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

package main

import (
	"io"

	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/spf13/cobra"
)

const (
	bumpDesc = `Bumps the semantic version within your git repository. The version bump 
is based on the conventional commit message from the last commit. Uplift
can also bump individual files in your repository to keep them in sync`
)

// bumps the version within a git repository tag, and bumps files

type bumpOptions struct {
	dryRun  bool
	verbose bool
}

func newBumpCmd(out io.Writer, cfg config.Uplift) *cobra.Command {
	opts := bumpOptions{}

	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Bump the semantic version",
		Long:  bumpDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run(out, cfg, opts)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&opts.dryRun, "dry-run", false, "take a practice bump")
	f.BoolVar(&opts.verbose, "verbose", false, "show me everything that happens")

	return cmd
}

func (o bumpOptions) Run(out io.Writer, cfg config.Uplift, opts bumpOptions) error {
	b := semver.NewBumper(out, semver.BumpOptions{
		Config:  cfg,
		DryRun:  opts.dryRun,
		Verbose: opts.verbose,
	})

	return b.Bump()
}
