package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Package struct {
	Pos          lexer.Position
	Name         string
	Functions    []Function `@@+`
	IsMain       bool
	FunctionsMap map[string]*Function
}

func (pkg *Package) ResolveNames(errorsCollector *errors.Collector) {
	pkg.FunctionsMap = make(map[string]*Function, len(pkg.Functions))

	for _, function := range pkg.Functions {
		_, exists := pkg.FunctionsMap[function.Name]
		if exists {
			errorsCollector.Err(
				function.Pos,
				"NameError",
				fmt.Sprintf("the function %s is redeclared", function.Name),
			)
		}

		pkg.FunctionsMap[function.Name] = &function

		function.ResolveNames(errorsCollector)
	}

	if pkg.IsMain {
		_, mainExists := pkg.FunctionsMap["main"]
		if !mainExists {
			errorsCollector.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is the main package and does not contain a main function", pkg.Name),
			)
		}
	} else {
		_, mainExists := pkg.FunctionsMap["main"]
		if mainExists {
			errorsCollector.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is not the main package but contains a main function", pkg.Name),
			)
		}
	}
}
