package alu

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
)

type LessThanUnsigned struct {
	tool.Node[LessThanUnsigned]
	typecheck.TypeCheck[struct {
		Type  types.Type
		Bytes int
	}]

	Left  memory.Read `parser:"'lt_u' '(' @@ ','"`
	Right memory.Read `parser:"@@ ')'"`
}

func (node *LessThanUnsigned) ResolveNames(ctx context.Context) error {
	if err := node.Left.ResolveNames(ctx); err != nil {
		return err
	}

	if err := node.Right.ResolveNames(ctx); err != nil {
		return err
	}

	return nil
}

func (node *LessThanUnsigned) ResolveTypes(target types.Type) error {
	left := node.Left.Type()
	right := node.Right.Type()

	if !left.Equals(right) {
		return node.Errorf("cannot compare %s and %s", left, right)
	}

	result := &types.ScalarType{Size: 1}

	if !target.CanBeAssigned(result) {
		return node.Errorf("cannot assign %s result to %s", result, target)
	}

	node.TypeCheckPass.Type = left

	return nil
}

func (node *LessThanUnsigned) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		*(*bool)(result) = *(*byte)(node.Left.Addr(vm)) < *(*byte)(node.Right.Addr(vm))
	case 2:
		*(*bool)(result) = *(*uint16)(node.Left.Addr(vm)) < *(*uint16)(node.Right.Addr(vm))
	case 4:
		*(*bool)(result) = *(*uint32)(node.Left.Addr(vm)) < *(*uint32)(node.Right.Addr(vm))
	case 8:
		*(*bool)(result) = *(*uint64)(node.Left.Addr(vm)) < *(*uint64)(node.Right.Addr(vm))
	}

	return nil
}
