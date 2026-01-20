package literal

import (
	"context"
	"math/bits"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Scalar struct {
	tool.Node[Scalar]
	typecheck.TypeCheck[struct {
		Type types.Type
	}]
	registeralloc.RegisterAllocation[struct {
		Target architecture.Register
	}]

	Value uint64 `parser:"@BitLiteral"`
}

func (node *Scalar) CalculateStatistics() {}

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

func (node *Scalar) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return nil, nil
}

func (node *Scalar) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Target = r
}

func (node *Scalar) GenerateVirtualMachineAssembly(p *assembly.P) error {
	p.AddLit(node.RegisterAllocationPass.Target.(bytecode.R), node.TypeCheckPass.Type.Bytes(), node.Value, node.Position())

	return nil
}
