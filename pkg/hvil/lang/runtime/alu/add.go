package alu

import (
	"errors"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
)

type Add struct {
	tool.Node[Add]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Left  memory.Read `parser:"'add' '(' @@ ','"`
	Right memory.Read `parser:"'@@ ')'"`
}

func (node *Add) ResolveNames(vars names.Scope[*stack.Decl], regs names.Scope[*memory.RegWrite]) error {
	return errors.Join(
		node.Left.ResolveNames(vars, regs),
		node.Right.ResolveNames(vars, regs),
	)
}

func (node *Add) ResolveTypes(target types.Type) error {
	if !target.CanBeAssigned(types.Void{}) {
		return node.Errorf("cannot assign void to %s", target)
	}

	//node.TypeCheckPass.Type = node.Param.Type()

	return nil
}

func (node *Add) Execute(vm *runtime.VirtualMachine, _ unsafe.Pointer) error {
	//var value any
	//switch node.TypeCheckPass.Type.Bytes() {
	//case 1:
	//	value = *(*byte)(node.Param.Addr(vm))
	//case 2:
	//	value = *(*uint16)(node.Param.Addr(vm))
	//case 4:
	//	value = *(*uint32)(node.Param.Addr(vm))
	//case 8:
	//	value = *(*uint64)(node.Param.Addr(vm))
	//}

	return nil
}
