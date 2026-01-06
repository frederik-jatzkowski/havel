package cmd

import "github.com/spf13/cobra"

func NewHvilCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hvil",
		Short: "Contains all commands for working with the Havel Intermediate Language (HVIL).",
	}

	cmd.AddCommand(NewHvilDumpCmd())
	cmd.AddCommand(NewHvilRunCmd())

	return cmd
}
