package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
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

func (pkg *Package) VisitLCR(visitor Visitor) {
	visitor.SetCurrentPackage(pkg)
	visitor.VisitPackage(pkg)

	for _, function := range pkg.Functions {
		function.VisitLCR(visitor)
	}
}
