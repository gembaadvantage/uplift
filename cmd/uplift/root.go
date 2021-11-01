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
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/spf13/cobra"
)

func newRootCmd(out io.Writer, args []string, ctx *context.Context) (*cobra.Command, error) {
	log.SetHandler(cli.Default)

	cmd := &cobra.Command{
		Use:          "uplift",
		Short:        "Semantic versioning the easy way",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if ctx.Debug {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	// Write persistent flags straight into the context
	pf := cmd.PersistentFlags()
	pf.BoolVar(&ctx.DryRun, "dry-run", false, "run without making any changes")
	pf.BoolVar(&ctx.Debug, "debug", false, "show me everything that happens")

	cmd.AddCommand(newVersionCmd(out),
		newBumpCmd(out, ctx),
		newCompletionCmd(out),
		newTagCmd(out, ctx),
		newReleaseCmd(out, ctx))

	return cmd, nil
}
