package cmd

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/frederik-jatzkowski/havel/internal/hvil/parser"
	"github.com/frederik-jatzkowski/havel/internal/hvil/pass"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
	"github.com/spf13/cobra"
)

//go:embed hvil_examine.html
var hvilExamineTemplate string

var examineCmd = &cobra.Command{
	Use:   "examine",
	Short: "Generates the ast for an hvil program and shows the results on an interactive web page.",
	Run: func(cmd *cobra.Command, args []string) {
		template := template.Must(template.New("root").Parse(hvilExamineTemplate))

		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			type Data struct {
				Output []string
				Ast    string
			}
			data := Data{}

			mainFilePath := args[0]
			mainFile, err := os.Open(mainFilePath)
			if err != nil {
				data.Output = append(data.Output, err.Error())
			}

			mainPkg, err := parser.Parse(mainFilePath, mainFile)
			if err != nil {
				data.Output = append(data.Output, err.Error())
			}

			mainPkg.Name = mainPkg.Pos.Filename
			mainPkg.IsMain = true

			program := parser.Program{
				Packages: []*parser.Package{
					&mainPkg,
				},
			}

			nameResolutionPass := pass.NameResolution{
				Result: errors.NewCollector(os.Stderr),
			}
			program.VisitCLR(&nameResolutionPass)

			if nameResolutionPass.Result.HasErrors() {
				for _, err := range nameResolutionPass.Result.Errors() {
					data.Output = append(data.Output, err.String())
				}
			}

			typeCheckPass := pass.TypeCheck{
				Result: errors.NewCollector(os.Stderr),
			}
			program.VisitCLR(&typeCheckPass)

			if typeCheckPass.Result.HasErrors() {
				for _, err := range typeCheckPass.Result.Errors() {
					data.Output = append(data.Output, err.String())
				}
			}

			astData, err := json.MarshalIndent(program, "", "  ")
			cobra.CheckErr(err)

			data.Ast = string(astData)

			template.Execute(w, data)
		}))

		err := http.ListenAndServe("localhost:8080", http.DefaultServeMux)
		cobra.CheckErr(err)
	},
}

func init() {
	RootCmd.AddCommand(examineCmd)
	RootCmd.Args = cobra.RangeArgs(1, 1)
}
