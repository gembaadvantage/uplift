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

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/spf13/cobra"
)

func newRootCmd(out io.Writer, args []string, ctx *context.Context) (*cobra.Command, error) {
	log.SetHandler(cli.Default)

	// support toggling of logging
	var silent bool

	cmd := &cobra.Command{
		Use:          "uplift",
		Short:        "Semantic versioning the easy way",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if ctx.Debug {
				log.SetLevel(log.InvalidLevel)
			}

			if silent {
				// Switch logging handler, to ensure all logging is discarded
				log.SetHandler(discard.Default)
			}
		},
	}

	// Write persistent flags straight into the context
	pf := cmd.PersistentFlags()
	pf.BoolVarP(&ctx.DryRun, "dry-run", "d", false, "run without making any changes")
	pf.BoolVarP(&ctx.Debug, "debug", "v", false, "show me everything that happens")
	pf.BoolVarP(&ctx.NoPush, "no-push", "n", false, "no changes will be pushed to the git remote")
	pf.BoolVarP(&silent, "silent", "s", false, "silence all logging")

	cmd.AddCommand(newVersionCmd(out),
		newBumpCmd(out, ctx),
		newCompletionCmd(out),
		newTagCmd(out, ctx),
		newReleaseCmd(out, ctx),
		newChangelogCmd(out, ctx))

	return cmd, nil
}
