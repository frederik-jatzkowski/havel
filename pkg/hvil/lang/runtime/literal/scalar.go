package literal

import (
	"math/bits"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/stack"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
)

type Scalar struct {
	tool.Node[Scalar]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Value uint64 `parser:"@BitLiteral"`
}

func (l *Scalar) ResolveNames(vars names.Scope[*stack.Decl], regs names.Scope[*memory.RegWrite]) (errs []error) {
	return nil
}

func (l *Scalar) ResolveTypes(target types.Type) (errs []error) {
	_, ok := target.(types.ScalarType)
	if !ok {
		return append(errs, l.Errorf("cannot assign scalar literal to %s", target))
	}

	requiredBytes := (bits.Len64(l.Value) + 7) / 8
	availableBytes := target.Bytes()
	if requiredBytes > availableBytes {
		errs = append(errs, l.Errorf("cannot assign scalar literal %d to %s: value too big", l.Value, target))
	}

	l.TypeCheckPass.Type = target

	return errs
}

func (l *Scalar) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	switch l.TypeCheckPass.Type.Bytes() {
	case 1:
		*(*byte)(result) = byte(l.Value)
	case 2:
		*(*uint16)(result) = uint16(l.Value)
	case 4:
		*(*uint32)(result) = uint32(l.Value)
	case 8:
		*(*uint64)(result) = l.Value
	}

	return nil
}
