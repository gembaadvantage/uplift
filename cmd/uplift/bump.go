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

	"github.com/spf13/cobra"
)

const (
	firstVersion = "v0.1.0"
)

type bumpOptions struct {
	first  string
	dryRun bool
}

func newBumpCmd(out io.Writer) *cobra.Command {
	opts := bumpOptions{}

	cmd := &cobra.Command{
		Use:   "bump",
		Short: "Bump the version of your application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run(out)
		},
	}

	f := cmd.Flags()
	f.StringVar(&opts.first, "first", firstVersion, "sets the first version of the initial bump")
	f.BoolVar(&opts.dryRun, "dry-run", false, "take a practice bump")

	return cmd
}

func (o bumpOptions) Run(out io.Writer) error {
	return nil
}
