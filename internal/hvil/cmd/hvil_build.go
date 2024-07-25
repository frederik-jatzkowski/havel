package cmd

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/frederik-jatzkowski/havel/internal/hvil/parser"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compiles a .hvil project.",
	Run: func(cmd *cobra.Command, args []string) {
		mainFilePath := args[0]
		mainFile, err := os.Open(mainFilePath)
		cobra.CheckErr(err)

		mainPkg, err := parser.Parse(mainFilePath, mainFile)
		cobra.CheckErr(err)

		mainPkg.Name = mainPkg.Pos.Filename
		mainPkg.IsMain = true

		program := parser.Program{
			Packages: []*parser.Package{
				&mainPkg,
			},
		}

		errorsCollector := errors.NewCollector(os.Stderr)

		program.GenerateBackLinks()
		program.ResolveNames(errorsCollector)

		if errorsCollector.HasErrors() {
			os.Exit(1)
		}

		data, err := json.MarshalIndent(program, "", "  ")
		cobra.CheckErr(err)

		removePos := regexp.MustCompile(`\n[\s]*"Pos": {\n.*\n.*\n.*\n.*\n.*},`)
		os.Stdout.Write(removePos.ReplaceAll(data, []byte{}))
		// os.Stdout.Write(data)

		// spew.Dump(mainPkg)
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	RootCmd.Args = cobra.RangeArgs(1, 1)
}
