package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Package struct {
	Name         string
	IsMain       bool
	Functions    []*Function `@@+`
	Pos          lexer.Position
	program      *Program
	functionsMap map[string]*Function
}

func (pkg *Package) GenerateBackLinks(program *Program) {
	pkg.program = program

	for _, function := range pkg.Functions {
		function.GenerateBackLinks(pkg)
	}
}

func (pkg *Package) ResolveNames(errorsCollector *errors.Collector) {
	pkg.functionsMap = make(map[string]*Function, len(pkg.Functions))
	for _, function := range pkg.Functions {
		_, exists := pkg.functionsMap[function.Name]
		if exists {
			errorsCollector.Err(
				function.Pos,
				"NameError",
				fmt.Sprintf("the function %s is redeclared in this package", function.Name),
			)
		}

		pkg.functionsMap[function.Name] = function
	}

	for _, function := range pkg.Functions {
		function.ResolveNames(errorsCollector)
	}

	if pkg.IsMain {
		_, mainExists := pkg.functionsMap["main"]
		if !mainExists {
			errorsCollector.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is the main package and does not contain a main function", pkg.Name),
			)
		}
	} else {
		_, mainExists := pkg.functionsMap["main"]
		if mainExists {
			errorsCollector.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is not the main package but contains a main function", pkg.Name),
			)
		}
	}
}
