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
	"errors"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/spf13/cobra"
)

func newRootCmd(args []string, ctx *context.Context) (*cobra.Command, error) {
	log.SetHandler(cli.Default)

	// Capture temporary flags
	var silent bool
	var pre string

	cmd := &cobra.Command{
		Use:          "uplift",
		Short:        "Semantic versioning the easy way",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if ctx.Debug {
				log.SetLevel(log.InvalidLevel)
			}

			if silent {
				// Switch logging handler, to ensure all logging is discarded
				log.SetHandler(discard.Default)
			}

			var err error
			ctx.Prerelease, ctx.Metadata, err = parsePrereleaseArg(pre)
			return err
		},
	}

	// Write persistent flags straight into the context
	pf := cmd.PersistentFlags()
	pf.BoolVar(&ctx.DryRun, "dry-run", false, "run without making any changes")
	pf.BoolVar(&ctx.Debug, "debug", false, "show me everything that happens")
	pf.BoolVar(&ctx.NoPush, "no-push", false, "no changes will be pushed to the git remote")
	pf.BoolVar(&silent, "silent", false, "silence all logging")
	pf.StringVar(&pre, "prerelease", "", "append a prerelease suffix to calculated semantic version")

	cmd.AddCommand(newVersionCmd(ctx),
		newBumpCmd(ctx),
		newCompletionCmd(ctx),
		newTagCmd(ctx),
		newReleaseCmd(ctx),
		newChangelogCmd(ctx))

	return cmd, nil
}

func parsePrereleaseArg(pre string) (string, string, error) {
	if pre == "" {
		return "", "", errors.New("prerelease suffix is blank")
	}

	// Has prefix been provided
	i := 0
	if pre[0] == '-' {
		i = 1
	}

	v, err := semver.Parse("1.0.0-" + pre[i:])
	if err != nil {
		return "", "", errors.New("invalid semantic versioning prerelease suffix")
	}

	return v.Prerelease, v.Metadata, nil
}
