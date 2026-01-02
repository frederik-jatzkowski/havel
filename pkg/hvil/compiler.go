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

func (c Compiler) Compile(path string, src io.Reader) (program.Program, error) {
	p, err := parser.Parse(path, src)
	if err != nil {
		return p, err
	}

	if err = p.ResolveNames(); err != nil {
		return p, err
	}

	if err = p.ResolveTypes(); err != nil {
		return p, err
	}

	if err = p.ResolveAddresses(); err != nil {
		return p, err
	}

	return p, nil
}
