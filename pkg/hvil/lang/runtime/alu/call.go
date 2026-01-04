package alu

import (
	"context"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
)

type Operation instruction.Operation

type Call struct {
	tool.Node[Call]

	Operation Operation `parser:"'alu' '.' @@"`
}

func (node *Call) ResolveNames(ctx context.Context) error {
	return node.Operation.ResolveNames(ctx)
}

func (node *Call) ResolveTypes(target types.Type) error {
	return node.Operation.ResolveTypes(target)
}

func (node *Call) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	return node.Operation.Execute(vm, result)
}
