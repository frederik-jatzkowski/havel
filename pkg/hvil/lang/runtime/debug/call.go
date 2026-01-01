package debug

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

	Op Op `parser:"'debug' '.' @@"`
}

func (c *Call) ResolveNames(vars names.Scope[*stack.Decl], regs names.Scope[*memory.RegWrite]) (errs []error) {
	return c.Op.ResolveNames(vars, regs)
}

func (c *Call) ResolveTypes(target types.Type) (errs []error) {
	return c.Op.ResolveTypes(target)
}

func (c *Call) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	return c.Op.Execute(vm, result)
}
