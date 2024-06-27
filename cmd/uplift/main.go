package main

import (
	"fmt"
	"os"
)

func main() {
	rootCmd := newRootCmd(os.Stdout)

	if err := rootCmd.Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
