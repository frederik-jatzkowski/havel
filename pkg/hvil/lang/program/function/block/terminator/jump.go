package terminator

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Jump struct {
	tool.Node[Jump]
	names.NameResolution[struct {
		Target *block.Block
	}]

	Target string `parser:"'goto':Keyword @Ident"`
}

var _ block.Terminator = (*Jump)(nil)

func (node *Jump) ResolveNames(
	_ names.Scope[*stack.Decl],
	_ names.Scope[*memory.RegWrite],
	blocks names.Scope[*block.Block],
) error {
	target, err := blocks.Find(node.Target)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Target = target

	return nil
}

func (node *Jump) ResolveTypes() error {
	return nil
}

func (node *Jump) Execute(vm *runtime.VirtualMachine) (*block.Block, error) {
	return node.NameResolutionPass.Target, nil
}
