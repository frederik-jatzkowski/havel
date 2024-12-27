package cmd

import "github.com/spf13/cobra"

var hvilCmd = &cobra.Command{
	Use:   "hvil",
	Short: "Contains all commands for working with the Havel Intermediate Language (HVIL).",
}

func init() {
	hvilCmd.Args = cobra.RangeArgs(1, 1)
}
