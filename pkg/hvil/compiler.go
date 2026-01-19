package hvil

import (
	"context"
	"io"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/parser"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine"
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

	if err = p.ResolveNames(context.Background()); err != nil {
		return p, err
	}

	if err = p.ResolveTypes(); err != nil {
		return p, err
	}

	if err = p.ResolveAddresses(virtualmachine.NewArchitecture()); err != nil {
		return p, err
	}

	return p, nil
}
