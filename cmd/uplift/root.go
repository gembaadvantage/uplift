package main

import (
	"io"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/spf13/cobra"
)

type globalOptions struct {
	DryRun                   bool
	Debug                    bool
	NoPush                   bool
	NoStage                  bool
	Silent                   bool
	IgnoreExistingPrerelease bool
	FilterOnPrerelease       bool
	IgnoreDetached           bool
	IgnoreShallow            bool
	ConfigDir                string
}
type rootCommand struct {
	Cmd  *cobra.Command
	Opts *globalOptions
}

func newRootCmd(out io.Writer) *rootCommand {
	log.SetHandler(cli.Default)

	rootCmd := &rootCommand{
		Opts: &globalOptions{},
	}

	cmd := &cobra.Command{
		Use:           "uplift",
		Short:         "Semantic versioning the easy way",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
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
	pf.BoolVar(&rootCmd.Opts.NoStage, "no-stage", false, "no changes will be git staged")
	pf.BoolVar(&rootCmd.Opts.Silent, "silent", false, "silence all logging")
	pf.BoolVar(&rootCmd.Opts.IgnoreDetached, "ignore-detached", false, "ignore reported git detached HEAD error")
	pf.BoolVar(&rootCmd.Opts.IgnoreShallow, "ignore-shallow", false, "ignore reported git shallow clone error")
	pf.BoolVar(&rootCmd.Opts.IgnoreExistingPrerelease, "ignore-existing-prerelease", false, "ignore any existing prerelease when calculating next semantic version")
	pf.BoolVar(&rootCmd.Opts.FilterOnPrerelease, "filter-on-prerelease", false, "filter tags that only match the provided prerelease when calculating next semantic version")

	cmd.AddCommand(newVersionCmd(out),
		newCompletionCmd(out),
		newBumpCmd(rootCmd.Opts, out).Cmd,
		newTagCmd(rootCmd.Opts, out).Cmd,
		newReleaseCmd(rootCmd.Opts, out).Cmd,
		newChangelogCmd(rootCmd.Opts, out).Cmd,
		newManPageCmd(out).Cmd,
		newCheckCmd(rootCmd.Opts, out),
	)

	rootCmd.Cmd = cmd
	return rootCmd
}
