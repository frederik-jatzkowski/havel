package alu

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Operation struct {
	tool.Node[Operation]
	tool.NotImplemented[Operation]

	Name string                 `parser:"'alu' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}

func (node *Operation) ResolveNames(vars names.Scope[*stack.Decl], regs names.Scope[*memory.RegWrite]) (errs []error) {
	return nil
}

func (node *Operation) ResolveTypes(target types.Type) (errs []error) {
	return append(errs, node.Errorf("not implemented"))
}

func (node *Operation) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	return node.Errorf("not implemented")
}
