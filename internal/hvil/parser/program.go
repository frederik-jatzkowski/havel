package parser

type Program struct {
	Packages   []*Package
	PackageMap map[string]*Package `parser:"" json:"-"`
}

func (program *Program) VisitLCR(visitor Visitor) {
	visitor.SetCurrentProgram(program)
	visitor.VisitProgram(program)

	for _, pkg := range program.Packages {
		pkg.VisitLCR(visitor)
	}
}
