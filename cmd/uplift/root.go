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
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/spf13/cobra"
)

type globalOpts struct {
	DryRun     bool
	Debug      bool
	NoPush     bool
	Silent     bool
	ConfigDir  string
	ConfigFile string
}

func newRootCmd(args []string) (*cobra.Command, error) {
	log.SetHandler(cli.Default)

	// Capture temporary flags
	var opts globalOpts
	var ctx *context.Context

	cmd := &cobra.Command{
		Use:          "uplift",
		Short:        "Semantic versioning the easy way",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(opts.ConfigDir)
			if err != nil {
				fmt.Printf("failed to load uplift config. %v", err)
				return err
			}
			ctx = context.New(cfg, os.Stdout)

			// Assign flags
			ctx.Debug = opts.Debug
			ctx.DryRun = opts.DryRun
			ctx.NoPush = opts.NoPush

			if ctx.Debug {
				log.SetLevel(log.InvalidLevel)
			}

			if opts.Silent {
				// Switch logging handler, to ensure all logging is discarded
				log.SetHandler(discard.Default)
			}

			return nil
		},
	}

	// Persistent flags to be written into the context
	pf := cmd.PersistentFlags()
	pf.StringVar(&opts.ConfigDir, "config-dir", currentWorkingDir, "a custom path to a directory containing uplift config")
	pf.BoolVar(&opts.DryRun, "dry-run", false, "run without making any changes")
	pf.BoolVar(&opts.Debug, "debug", false, "show me everything that happens")
	pf.BoolVar(&opts.NoPush, "no-push", false, "no changes will be pushed to the git remote")
	pf.BoolVar(&opts.Silent, "silent", false, "silence all logging")

	cmd.AddCommand(newVersionCmd(ctx),
		newBumpCmd(ctx),
		newCompletionCmd(ctx),
		newTagCmd(ctx),
		newReleaseCmd(ctx),
		newChangelogCmd(ctx))

	return cmd, nil
}
