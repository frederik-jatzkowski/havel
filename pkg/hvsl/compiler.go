package hvsl

import (
	"context"
	"io"

	"github.com/frederik-jatzkowski/havel/pkg/hvsl/internal/pass/parser"
	"github.com/frederik-jatzkowski/havel/pkg/hvsl/program"
)

type Compiler struct {
}

func NewCompiler() Compiler {
	return Compiler{}
}

func (c Compiler) Compile(path string, src io.Reader) (*program.Program, error) {
	p, err := parser.Parse(path, src)
	if err != nil {
		return &p, err
	}

	ctx := context.Background()

	if err = p.ResolveNames(ctx); err != nil {
		return &p, err
	}

	if err = p.ResolveTypes(ctx); err != nil {
		return &p, err
	}

	return &p, nil
}
