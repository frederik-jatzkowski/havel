package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
)

type SizeOf struct {
	tool.Node[SizeOf]
	typecheck.TypeCheck[struct {
		Result types.Type
		IsVoid bool
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Param types.Type `parser:"'sizeof' '(' @@ ')'"`
}

func (node *SizeOf) ResolveNames(ctx context.Context) error {
	return nil
}

func (node *SizeOf) ResolveTypes(target types.Type) error {
	switch target.(type) {
	case *types.Scalar:
	case *types.Void:
		node.TypeCheckPass.IsVoid = true
	default:
		return node.Errorf("cannot assign size to %s", target)
	}

	node.TypeCheckPass.Result = target

	return nil
}

func (node *SizeOf) CalculateStatistics(ctx context.Context) {}

func (node *SizeOf) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return nil, nil
}

func (node *SizeOf) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *SizeOf) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if !node.TypeCheckPass.IsVoid {
		p.AddLit(node.RegisterAllocationPass.Result.(bytecode.R), node.TypeCheckPass.Result.Bytes(), uint64(node.Param.Bytes()), node.Position())
	}

	return nil
}
