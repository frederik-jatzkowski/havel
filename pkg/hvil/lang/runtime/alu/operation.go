package alu

import (
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"unsafe"
)

type Operation struct {
	tool.Node[Operation]
	tool.NotImplemented[Operation]

	Name string                 `parser:"'alu' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}

func (o Operation) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return nil
}

func (o Operation) ResolveTypes(target types.Type) (errs []error) {
	return append(errs, o.Errorf("not implemented"))
}

func (o Operation) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	return o.Errorf("not implemented")
}
