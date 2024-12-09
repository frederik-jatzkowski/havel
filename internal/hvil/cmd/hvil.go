package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "hvil",
	Short: "Contains all commands for working with the Havel Intermediate Language (HVIL).",
}

func init() {
	RootCmd.Args = cobra.RangeArgs(1, 1)
}
