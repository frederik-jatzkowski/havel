package alu

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Operation instruction.Operation

type Call struct {
	tool.Node[Call]

	Operation Operation `parser:"'alu' '.' @@"`
}

func (node *Call) ResolveNames(vars names.Scope[*stack.Decl], regs names.Scope[*memory.RegWrite]) error {
	return node.Operation.ResolveNames(vars, regs)
}

func (node *Call) ResolveTypes(target types.Type) error {
	return node.Operation.ResolveTypes(target)
}

func (node *Call) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	return node.Operation.Execute(vm, result)
}
