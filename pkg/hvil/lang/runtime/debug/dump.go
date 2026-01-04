package debug

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
)

type Dump struct {
	tool.Node[Dump]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Param memory.Read `parser:"'dump' '(' @@ ')'"`
}

func (node *Dump) ResolveNames(ctx context.Context) error {
	return node.Param.ResolveNames(ctx)
}

func (node *Dump) ResolveTypes(target types.Type) error {
	if !target.CanBeAssigned(types.Void{}) {
		return node.Errorf("cannot assign void to %s", target)
	}

	node.TypeCheckPass.Type = node.Param.Type()

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

	memKind := "unknown"
	switch node.Param.(type) {
	case *memory.RegRead:
		memKind = "register"
	case *memory.VarRead:
		memKind = "variable"
	}

	metaString := fmt.Sprintf("(%s %s '%s')", node.Param.Type(), memKind, node.Param.Identifier())

	_, err := fmt.Fprintln(vm.Stdout, node.Pos, value, metaString)

	return err
}
