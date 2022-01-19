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
	"github.com/spf13/cobra"
)

type globalOptions struct {
	DryRun    bool
	Debug     bool
	NoPush    bool
	Silent    bool
	ConfigDir string
}
type rootCommand struct {
	Cmd  *cobra.Command
	Opts globalOptions
}

func newRootCmd(out io.Writer) *rootCommand {
	log.SetHandler(cli.Default)

	rootCmd := &rootCommand{}

	cmd := &cobra.Command{
		Use:          "uplift",
		Short:        "Semantic versioning the easy way",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if rootCmd.Opts.Debug {
				log.SetLevel(log.InvalidLevel)
			}

			if rootCmd.Opts.Silent {
				// Switch logging handler, to ensure all logging is discarded
				log.SetHandler(discard.Default)
			}
		},
	}

	// Persistent flags to be written into the context
	pf := cmd.PersistentFlags()
	pf.StringVar(&rootCmd.Opts.ConfigDir, "config-dir", currentWorkingDir, "a custom path to a directory containing uplift config")
	pf.BoolVar(&rootCmd.Opts.DryRun, "dry-run", false, "run without making any changes")
	pf.BoolVar(&rootCmd.Opts.Debug, "debug", false, "show me everything that happens")
	pf.BoolVar(&rootCmd.Opts.NoPush, "no-push", false, "no changes will be pushed to the git remote")
	pf.BoolVar(&rootCmd.Opts.Silent, "silent", false, "silence all logging")

	cmd.AddCommand(newVersionCmd(out),
		newCompletionCmd(out),
		newBumpCmd(rootCmd.Opts, out).Cmd,
		newTagCmd(rootCmd.Opts, out).Cmd,
		newReleaseCmd(rootCmd.Opts, out).Cmd,
		newChangelogCmd(rootCmd.Opts, out).Cmd)

	rootCmd.Cmd = cmd
	return rootCmd
}
