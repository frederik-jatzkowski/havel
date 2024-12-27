package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/parser"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Generates the ast for an hvil program and dumps the result into a json file.",
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]

		err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(info.Name(), ".hvil") {
				return nil
			}

			file, err := os.Open(path)
			cobra.CheckErr(err)
			defer file.Close()

			program, err := parser.Parse(path, file)
			if err != nil {
				fmt.Println(err)

				return nil
			}

			errs := program.ResolveNames()
			if len(errs) > 0 {
				for _, err := range errs {
					fmt.Println(err)
				}

				return nil
			}

			data, err := json.MarshalIndent(program, "", "  ")
			cobra.CheckErr(err)

			astFile, err := os.Create(fmt.Sprintf("%s.ast.json", path))
			cobra.CheckErr(err)
			defer astFile.Close()

			_, err = astFile.Write(data)
			cobra.CheckErr(err)

			return nil
		})
		cobra.CheckErr(err)
	},
}

func init() {
	hvilCmd.AddCommand(dumpCmd)
}
