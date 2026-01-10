package terminator

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Return struct {
	tool.Node[Return]

	Token string `parser:"@'return':Keyword" json:"-"`
}

var _ block.Terminator = (*Return)(nil)

func (node *Return) ResolveNames(_ context.Context) error {
	return nil
}

func (node *Return) ResolveTypes() error {
	return nil
}

func (node *Return) GenerateVirtualMachineAssembly(p *assembly.P, isMain bool) error {
	if isMain {
		// exit code 0
		tmp := bytecode.R(255)
		p.AddLit(tmp, 1, 0, node.Position())
		p.AddI1R(bytecode.OPExit, tmp, node.Position())
	} else {
		//TODO implement me
		panic("implement me")
	}

	return nil
}

func (node *Return) Execute(vm *runtime.VirtualMachine) (*block.Block, error) {
	// there is no next block after a return statement
	return nil, nil
}
