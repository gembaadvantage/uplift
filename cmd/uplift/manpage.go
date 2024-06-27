package main

import (
	"fmt"
	"io"

	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

type manPageCommand struct {
	Cmd *cobra.Command
}

func newManPageCmd(out io.Writer) *manPageCommand {
	manCmd := &manPageCommand{}

	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generate man pages for uplift",
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			mp, err := mcobra.NewManPage(1, manCmd.Cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(out, mp.Build(roff.NewDocument()))
			return err
		},
	}

	manCmd.Cmd = cmd
	return manCmd
}
