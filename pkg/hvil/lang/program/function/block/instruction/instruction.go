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
	Operation Op           `parser:"@@ ';'"`
}

func (i *Instruction) ResolveNames(
	vars names.Scope[*stack.Decl],
	regs names.Scope[*memory.RegWrite],
) (errs []error) {
	if i.Result != nil {
		errs = append(errs, i.Result.ResolveNames(vars, regs)...)
	}

	errs = append(errs, i.Operation.ResolveNames(vars, regs)...)

	return errs
}

func (i *Instruction) ResolveTypes() (errs []error) {
	if i.Result != nil {
		return i.Operation.ResolveTypes(i.Result.Type())
	}

	return i.Operation.ResolveTypes(types.Void{})
}

func (i *Instruction) Execute(vm *runtime.VirtualMachine) error {
	var result unsafe.Pointer
	if i.Result != nil {
		result = i.Result.Addr(vm)
	}

	return i.Operation.Execute(vm, result)
}
