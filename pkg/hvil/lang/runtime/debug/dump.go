package debug

import (
	"fmt"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"unsafe"
)

type Dump struct {
	tool.Node[Dump]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Param memory.Read `parser:"'dump' '(' @@ ')'"`
}

func (d *Dump) ResolveNames(vars names.Scope[memory.VarDecl], regs names.Scope[memory.RegWrite]) (errs []error) {
	return d.Param.ResolveNames(vars, regs)
}

func (d *Dump) ResolveTypes(target types.Type) (errs []error) {
	if !target.CanBeAssigned(types.Void{}) {
		errs = append(errs, d.Errorf("cannot assign void to %s", target))
	}

	d.TypeCheckPass.Type = d.Param.Type()

	return errs
}

func (d *Dump) Execute(vm *runtime.VirtualMachine, _ unsafe.Pointer) (err error) {

	var value any
	switch d.TypeCheckPass.Type.Bytes() {
	case 1:
		value = *(*byte)(d.Param.Addr(vm))
	case 2:
		value = *(*uint16)(d.Param.Addr(vm))
	case 4:
		value = *(*uint32)(d.Param.Addr(vm))
	case 8:
		value = *(*uint64)(d.Param.Addr(vm))
	}

	memKind := "register"
	if _, ok := d.Param.(*memory.VarRead); ok {
		memKind = "variable"
	}

	metaString := fmt.Sprintf("(%s %s '%s')", d.Param.Type(), memKind, d.Param.Identifier())

	_, err = fmt.Fprintln(vm.Stdout, d.Pos, value, metaString)

	return err
}
