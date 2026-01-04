package literal

import (
	"context"
	"math/bits"
	"unsafe"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
)

type Scalar struct {
	tool.Node[Scalar]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]

	Value uint64 `parser:"@BitLiteral"`
}

func (node *Scalar) ResolveNames(_ context.Context) error {
	return nil
}

func (node *Scalar) ResolveTypes(target types.Type) error {
	_, ok := target.(*types.ScalarType)
	if !ok {
		return node.Errorf("cannot assign scalar literal to %s", target)
	}

	requiredBytes := (bits.Len64(node.Value) + 7) / 8
	availableBytes := target.Bytes()
	if requiredBytes > availableBytes {
		return node.Errorf("cannot assign scalar literal %d to %s: value too big", node.Value, target)
	}

	node.TypeCheckPass.Type = target

	return nil
}

func (node *Scalar) Execute(vm *runtime.VirtualMachine, result unsafe.Pointer) error {
	switch node.TypeCheckPass.Type.Bytes() {
	case 1:
		*(*byte)(result) = byte(node.Value)
	case 2:
		*(*uint16)(result) = uint16(node.Value)
	case 4:
		*(*uint32)(result) = uint32(node.Value)
	case 8:
		*(*uint64)(result) = node.Value
	}

	return nil
}
