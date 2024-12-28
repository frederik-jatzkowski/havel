package cmd

import (
	"errors"
	"github.com/frederik-jatzkowski/havel/pkg/hvil"
	"github.com/spf13/cobra"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Compiles a HVIL file and executes it.",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		compiler := hvil.NewCompiler()

		_, errs := compiler.Compile(filePath, file)
		if len(errs) > 0 {
			err = errors.New("compilation failed")
			for _, err2 := range errs {
				err = errors.Join(err, err2)
			}
			return err
		}

		return nil
	},
}

func init() {
	hvilCmd.AddCommand(runCmd)
}
