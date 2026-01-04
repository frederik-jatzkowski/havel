package local

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Call struct {
	tool.Node[Call]
	names.NameResolution[struct {
		Decl *function.Function
	}]
	tool.NotImplemented[Call]

	Name string                 `parser:"'local' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}

func (node *Call) ResolveNames(vars names.Scope[*stack.Decl], regs names.Scope[*memory.RegWrite]) error {
	return nil
}

func (node *Call) ResolveTypes(target types.Type) error {
	return nil
}

func (node *Call) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	return nil
}
