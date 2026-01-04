package terminator

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Return struct {
	tool.Node[Return]

	Token string `parser:"@'return':Keyword" json:"-"`
}

var _ block.Terminator = (*Return)(nil)

func (node *Return) ResolveNames(
	_ names.Scope[*stack.Decl],
	_ names.Scope[*memory.RegWrite],
	_ names.Scope[*block.Block],
) error {
	return nil
}

func (node *Return) Execute(vm *runtime.VirtualMachine) (*block.Block, error) {
	// there is no next block after a return statement
	return nil, nil
}
