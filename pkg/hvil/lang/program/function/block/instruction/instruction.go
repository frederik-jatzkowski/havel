package instruction

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Instruction struct {
	tool.Node[Instruction]

	Result    memory.Write `parser:"(@@ '=')?"`
	Operation Operation    `parser:"@@ ';'"`
}

func (node *Instruction) ResolveNames(
	vars names.Scope[*stack.Decl],
	regs names.Scope[*memory.RegWrite],
) error {
	if node.Result != nil {
		if err := node.Result.ResolveNames(vars, regs); err != nil {
			return err
		}
	}

	return node.Operation.ResolveNames(vars, regs)
}

func (node *Instruction) ResolveTypes() error {
	if node.Result != nil {
		return node.Operation.ResolveTypes(node.Result.Type())
	}

	return node.Operation.ResolveTypes(types.Void{})
}

func (node *Instruction) Execute(vm *runtime.VirtualMachine) error {
	var result unsafe.Pointer
	if node.Result != nil {
		result = node.Result.Addr(vm)
	}

	return node.Operation.Execute(vm, result)
}
