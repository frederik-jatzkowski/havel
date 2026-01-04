package terminator

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Conditional struct {
	tool.Node[Conditional]
	names.NameResolution[struct {
		Then, Else *block.Block
	}]

	Condition memory.Read `parser:"'if':Keyword @@"`
	Then      string      `parser:"'then':Keyword @Ident"`
	Else      string      `parser:"'else':Keyword @Ident"`
}

var _ block.Terminator = (*Conditional)(nil)

func (node *Conditional) ResolveNames(
	vars names.Scope[*stack.Decl],
	regs names.Scope[*memory.RegWrite],
	blocks names.Scope[*block.Block],
) error {
	thenTarget, err := blocks.Find(node.Then)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Then = thenTarget

	elseTarget, err := blocks.Find(node.Else)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Else = elseTarget

	return node.Condition.ResolveNames(vars, regs)
}

func (node *Conditional) ResolveTypes() error {
	actualType := node.Condition.Type()

	if !actualType.Equals(&types.ScalarType{Size: 1}) {
		return node.Errorf("condition must be 1 byte but was %s", actualType)
	}

	return nil
}

func (node *Conditional) Execute(vm *runtime.VirtualMachine) (*block.Block, error) {
	cond := *(*byte)(node.Condition.Addr(vm))
	if cond > 0 {
		return node.NameResolutionPass.Then, nil
	}

	return node.NameResolutionPass.Else, nil
}
