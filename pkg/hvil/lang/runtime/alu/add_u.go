package alu

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
)

type AddUnsigned struct {
	tool.Node[AddUnsigned]
	typecheck.TypeCheck[struct {
		Type  types.Type
		Bytes int
	}]

	Left  memory.Read `parser:"'add_u' '(' @@ ','"`
	Right memory.Read `parser:"@@ ')'"`
}

func (node *AddUnsigned) ResolveNames(vars names.Scope[*stack.Decl], regs names.Scope[*memory.RegWrite]) error {
	if err := node.Left.ResolveNames(vars, regs); err != nil {
		return err
	}

	if err := node.Right.ResolveNames(vars, regs); err != nil {
		return err
	}

	return nil
}

func (node *AddUnsigned) ResolveTypes(target types.Type) error {
	left := node.Left.Type()
	right := node.Right.Type()

	if !left.Equals(right) {
		return node.Errorf("cannot add %s and %s", left, right)
	}

	if !target.CanBeAssigned(left) {
		return node.Errorf("cannot assign %s result to %s", left, target)
	}

	node.TypeCheckPass.Type = left

	return nil
}

func (node *AddUnsigned) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		*(*byte)(result) = *(*byte)(node.Left.Addr(vm)) + *(*byte)(node.Right.Addr(vm))
	case 2:
		*(*uint16)(result) = *(*uint16)(node.Left.Addr(vm)) + *(*uint16)(node.Right.Addr(vm))
	case 4:
		*(*uint32)(result) = *(*uint32)(node.Left.Addr(vm)) + *(*uint32)(node.Right.Addr(vm))
	case 8:
		*(*uint64)(result) = *(*uint64)(node.Left.Addr(vm)) + *(*uint64)(node.Right.Addr(vm))
	}

	return nil
}
