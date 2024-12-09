package program

import (
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/program/function"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/lang/tool"
	"github.com/frederik-jatzkowski/havel/internal/hvil/compiler/pass"
)

type Program struct {
	tool.Node[Program]
	pass.NameResolution[struct {
		Functions map[string]*function.Function
	}]

	Functions []function.Function `parser:"@@+"`
}

func (p *Program) ResolveNames() (errs []error) {
	p.NameResolutionPass.Functions = make(map[string]*function.Function)

	for i := 0; i < len(p.Functions); i++ {
		f := &p.Functions[i]

		_, exists := p.NameResolutionPass.Functions[f.Name]
		if exists {
			errs = append(errs, f.Errorf("function '%s' is already defined", f.Name))
		}

		p.NameResolutionPass.Functions[f.Name] = f
	}

	return errs
}
