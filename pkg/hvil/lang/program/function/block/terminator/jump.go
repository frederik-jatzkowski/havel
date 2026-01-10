package terminator

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
)

type Jump struct {
	tool.Node[Jump]
	names.NameResolution[struct {
		Target *block.Block
	}]

	Target string `parser:"'goto':Keyword @Ident"`
}

var _ block.Terminator = (*Jump)(nil)

func (node *Jump) ResolveNames(ctx context.Context) error {
	target, err := block.FromCtx(ctx, node.Target)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Target = target

	return nil
}

func (node *Jump) ResolveTypes() error {
	return nil
}

func (node *Jump) GenerateVirtualMachineAssembly(p *assembly.P, isMain bool) error {
	//TODO implement me
	panic("implement me")
}

func (node *Jump) Execute(vm *runtime.VirtualMachine) (*block.Block, error) {
	return node.NameResolutionPass.Target, nil
}
