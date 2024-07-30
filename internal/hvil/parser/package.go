package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Package struct {
	Name         string
	IsMain       bool
	Functions    []*Function `parser:"@@+"`
	Pos          lexer.Position
	FunctionsMap map[string]*Function `parser:"" json:"-"`
}

func (pkg *Package) VisitLCR(visitor Visitor) {
	visitor.SetCurrentPackage(pkg)
	visitor.VisitPackage(pkg)

	for _, function := range pkg.Functions {
		function.VisitLCR(visitor)
	}
}
