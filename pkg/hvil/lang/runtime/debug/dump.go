package debug

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Dump struct {
	tool.Node[Dump]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Param instruction.MemoryRead `parser:"'dump' '(' @@ ')'"`
}

func (node *Dump) ResolveNames(ctx context.Context) error {
	return node.Param.ResolveNames(ctx)
}

func (node *Dump) ResolveTypes(target types.Type) error {
	if !target.CanBeAssigned(&types.Void{}) {
		return node.Errorf("cannot assign void to %s", target)
	}

	node.TypeCheckPass.Type = node.Param.Type()

	return nil
}

func (node *Dump) CalculateStatistics() {}

func (node *Dump) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return node.Param.AllocateRegisters(scope)
}

func (node *Dump) SetResultRegister(r architecture.Register) {
	panic(fmt.Sprintf("target register assigned to %T, which returns void", node))
}

func (node *Dump) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Param.GenerateVirtualMachineAssembly(p); err != nil {
		return err
	}

	p.AddI1R(bytecode.OPDebugDump, node.Param.Register().(bytecode.R), node.Position())

	return nil
}

func (node *Dump) Execute(vm *runtime.VirtualMachine, _ unsafe.Pointer) error {
	var value any
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		value = *(*byte)(node.Param.Addr(vm))
	case 2:
		value = *(*uint16)(node.Param.Addr(vm))
	case 4:
		value = *(*uint32)(node.Param.Addr(vm))
	case 8:
		value = *(*uint64)(node.Param.Addr(vm))
	}

	_, err := fmt.Fprintf(vm.Stdout, "%s register content: %d\n", node.Position(), value)

	return err
}
