/*
Copyright (c) 2022 Gemba Advantage

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
