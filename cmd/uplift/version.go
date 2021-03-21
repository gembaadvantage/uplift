package main

import (
	"fmt"
	"io"

	"github.com/gembaadvantage/uplift/internal/version"
	"github.com/spf13/cobra"
)

type versionOptions struct {
	short bool
}

func newVersionCmd(out io.Writer) *cobra.Command {
	opts := versionOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the build time version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.run(out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&opts.short, "short", false, "only print the semantic version number")

	return cmd
}

func (o versionOptions) run(out io.Writer) error {
	fmt.Fprintln(out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	if short {
		return version.Short()
	}

	return fmt.Sprintf("%#v", version.Long())
}
