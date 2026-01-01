package hvil

import (
	"io"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/parser"
)

type Compiler struct {
}

func NewCompiler() Compiler {
	return Compiler{}
}

func (c Compiler) Compile(path string, src io.Reader) (program.Program, []error) {
	p, err := parser.Parse(path, src)
	if err != nil {
		return p, []error{err}
	}

	errs := p.ResolveNames()
	if len(errs) > 0 {
		return p, errs
	}

	errs = p.ResolveTypes()
	if len(errs) > 0 {
		return p, errs
	}

	errs = p.ResolveAddresses()
	if len(errs) > 0 {
		return p, errs
	}

	return p, nil
}
