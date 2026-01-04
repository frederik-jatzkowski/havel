package terminator

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Conditional struct {
	tool.Node[Conditional]
	tool.NotImplemented[Conditional]

	Condition memory.Read `parser:"'if':Keyword @@"`
	True      string      `parser:"'then':Keyword @Ident"`
	False     string      `parser:"'else':Keyword @Ident"`
}

var _ block.Terminator = (*Conditional)(nil)

func (node *Conditional) ResolveNames(
	vars names.Scope[*stack.Decl],
	regs names.Scope[*memory.RegWrite],
	blocks names.Scope[*block.Block],
) error {
	//TODO implement me
	panic("implement me")
}

func (node *Conditional) Execute(vm *runtime.VirtualMachine) (*block.Block, error) {
	//TODO implement me
	panic("implement me")
}
