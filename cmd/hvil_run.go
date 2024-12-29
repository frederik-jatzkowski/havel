package cmd

import (
	"errors"
	"github.com/frederik-jatzkowski/havel/pkg/hvil"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
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

		program, errs := compiler.Compile(filePath, file)
		if len(errs) > 0 {
			err = errors.New("compilation failed")
			for _, err2 := range errs {
				err = errors.Join(err, err2)
			}
			return err
		}

		vm := runtime.New(
			1024,
			os.Stdin,
			os.Stdout,
			os.Stderr,
		)

		err = program.Execute(vm)
		if err != nil {
			return errors.Join(
				errors.New("runtime error"),
				err,
			)
		}

		return nil
	},
}

func init() {
	hvilCmd.AddCommand(runCmd)
}
