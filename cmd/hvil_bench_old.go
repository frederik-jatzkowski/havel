package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/frederik-jatzkowski/havel/pkg/hvil"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
)

func NewHvilBenchOldCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bench-old",
		Short: "Compiles a HVIL file and executes it n times. Averages the runtimes.",
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
				io.Discard,
				io.Discard,
			)

			n, err := cmd.Flags().GetInt("executions")
			cobra.CheckErr(err)

			start := time.Now()

			for range n {
				err = program.Execute(vm)
				if err != nil {
					return errors.Join(
						errors.New("runtime error"),
						err,
					)
				}
			}

			fmt.Printf(
				"\nprogram execution took %d ns (%d executions)\n",
				time.Since(start).Nanoseconds()/int64(n),
				n,
			)

			return nil
		},
	}

	cmd.Flags().IntP("executions", "n", 1000, "Number of executions for benchmarking.")

	return cmd
}
