package parser

import (
	"fmt"

	"github.com/frederik-jatzkowski/havel/internal/tooling/errors"
)

type Program struct {
	Packages   []*Package
	packageMap map[string]*Package
}

func (program *Program) GenerateBackLinks() {
	for _, pkg := range program.Packages {
		pkg.GenerateBackLinks(program)
	}
}

func (program *Program) ResolveNames(errorsCollector *errors.Collector) {
	program.packageMap = make(map[string]*Package, len(program.Packages))
	for _, pkg := range program.Packages {
		_, exists := program.packageMap[pkg.Name]
		if exists {
			errorsCollector.Err(
				pkg.Pos,
				"NameError",
				fmt.Sprintf("the package %s is redeclared in this program", pkg.Name),
			)
		}

		program.packageMap[pkg.Name] = pkg
	}

	for _, pkg := range program.Packages {
		pkg.ResolveNames(errorsCollector)
	}
}
