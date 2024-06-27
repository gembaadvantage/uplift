package main

import (
	"io"

	"github.com/spf13/cobra"
)

func newCheckCmd(gopts *globalOptions, _ io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if a configuration file is valid",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := loadConfig(gopts.ConfigDir)
			if err != nil {
				return err
			}

			return cfg.Validate()
		},
	}

	return cmd
}
