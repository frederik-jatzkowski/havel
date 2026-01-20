package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/frederik-jatzkowski/havel/pkg/hvil"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine"
)

func NewHvilRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Compiles a HVIL file and executes it.",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			file, err := os.Open(filePath)
			cobra.CheckErr(err)

			compiler := hvil.NewCompiler()

			program, err := compiler.Compile(filePath, file)
			cobra.CheckErr(err)

			asm, err := program.GenerateVirtualMachineAssembly()
			cobra.CheckErr(err)

			byteCode, err := asm.Assemble()
			cobra.CheckErr(err)

			vm := virtualmachine.New(1024*1024, os.Stdin, os.Stdout, os.Stderr)

			return vm.Execute(byteCode)
		},
	}

	return cmd
}
