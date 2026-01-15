package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/frederik-jatzkowski/havel/pkg/hvil"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
)

func NewHvilRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Compiles a HVIL file and executes it.",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}

			compiler := hvil.NewCompiler()

			program, err := compiler.Compile(filePath, file)
			if err != nil {
				return fmt.Errorf("compilation failed:\n %w", err)
			}

			vm := runtime.New(
				1024*1024,
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

	return cmd
}
