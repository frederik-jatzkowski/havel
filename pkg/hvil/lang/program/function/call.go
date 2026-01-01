package function

import (
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
)

type Call struct {
	tool.Node[Call]
	names.NameResolution[struct {
		Decl *Function
	}]
	tool.NotImplemented[Call]

	Name string                 `parser:"'local' '.' @Ident"`
	Args tool.List[memory.Read] `parser:"'(' @@ ')'"`
}

func (node *Call) ResolveNames(vars names.Scope[*stack.Decl], regs names.Scope[*memory.RegWrite]) (errs []error) {
	return nil
}

func (node *Call) ResolveTypes(target types.Type) (errs []error) {
	return append(errs, node.Errorf("not implemented"))
}

func (node *Call) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	return node.Errorf("not implemented")
}
