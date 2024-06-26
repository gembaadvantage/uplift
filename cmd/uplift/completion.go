package main

import (
	"io"

	"github.com/spf13/cobra"
)

const (
	bashDesc = `Generate an uplift completion script for the bash shell. 
To use bash completions ensure you have them installed and enabled

To load completions in your current shell session:
  $ source <(uplift completion bash)

To load completions for every new session, execute once:
## Linux:
  $ uplift completion bash > /etc/bash_completion.d/uplift

## MacOS:
  $ uplift completion bash > /usr/local/etc/bash_completion.d/uplift`

	zshDesc = `Generate an uplift completion script for the zsh shell.

To load completions in your current shell session:
  $ source <(uplift completion zsh)

To load completions for every new session, execute once:
  $ uplift completion zsh > "${fpath[1]}/_uplift"

Alternatively install the uplift plugin with oh-my-zsh`

	fishDesc = `Generate uplift completion script for the fish shell.

To load completions in your current shell session:
  $ uplift completion fish | source

To load completions for every new session, execute once:
  $ uplift completion fish > ~/.config/fish/completions/uplift.fish

** You will need to start a new shell for this setup to take effect. **`

	noDescFlag     = "no-descriptions"
	noDescFlagText = "disable completion descriptions"
)

type completionOptions struct {
	noDescriptions bool
	shell          string
}

func newCompletionCmd(out io.Writer) *cobra.Command {
	opts := completionOptions{}

	cmd := &cobra.Command{
		Use:                   "completion",
		Short:                 "Generate completion script for your target shell",
		Long:                  "Generate an uplift completion script for either the bash, zsh or fish shells",
		DisableFlagsInUseLine: true,
	}

	bash := &cobra.Command{
		Use:                   "bash",
		Short:                 "generate a bash shell completion script",
		Long:                  bashDesc,
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.shell = "bash"
			return opts.Run(out, cmd)
		},
	}

	zsh := &cobra.Command{
		Use:   "zsh",
		Short: "generate a zsh shell completion script",
		Long:  zshDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.shell = "zsh"
			return opts.Run(out, cmd)
		},
	}
	zsh.Flags().BoolVar(&opts.noDescriptions, noDescFlag, false, noDescFlagText)

	fish := &cobra.Command{
		Use:   "fish",
		Short: "generate a fish shell completion script",
		Long:  fishDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.shell = "fish"
			return opts.Run(out, cmd)
		},
	}
	fish.Flags().BoolVar(&opts.noDescriptions, noDescFlag, false, noDescFlagText)

	cmd.AddCommand(bash, zsh, fish)
	return cmd
}

func (o completionOptions) Run(out io.Writer, cmd *cobra.Command) error {
	var err error
	switch o.shell {
	case "bash":
		err = cmd.Root().GenBashCompletionV2(out, !o.noDescriptions)
	case "zsh":
		if o.noDescriptions {
			err = cmd.Root().GenZshCompletionNoDesc(out)
		} else {
			err = cmd.Root().GenZshCompletion(out)
		}
	case "fish":
		err = cmd.Root().GenFishCompletion(out, !o.noDescriptions)
	}

	return err
}
