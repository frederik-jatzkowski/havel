package cmd

import (
	"encoding/json"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass/parser"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compiles a .hvil project.",
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		file, err := os.Open(filePath)
		cobra.CheckErr(err)

		program, err := parser.Parse(filePath, file)
		cobra.CheckErr(err)

		data, err := json.MarshalIndent(program, "", "  ")
		cobra.CheckErr(err)

		removePos := regexp.MustCompile(`\n[\s]*"Pos": {\n.*\n.*\n.*\n.*\n.*},`)
		os.Stdout.Write(removePos.ReplaceAll(data, []byte{}))
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
