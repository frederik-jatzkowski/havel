package cmd

import (
	"fmt"
	"os"

	hvilCmd "github.com/frederik-jatzkowski/havel/internal/hvil/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "havel",
}

func init() {
	rootCmd.AddCommand(hvilCmd.RootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
