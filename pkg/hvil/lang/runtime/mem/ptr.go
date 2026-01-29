package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/memory"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Ptr struct {
	tool.Node[Ptr]
	typecheck.TypeCheck[struct {
		Type   types.Type
		IsVoid bool
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Var *memory.VarRead `parser:"'ptr' '(' @@ ')'"`
}

func (node *Ptr) ResolveNames(ctx context.Context) error {
	return node.Var.ResolveNames(ctx)
}

func (node *Ptr) ResolveTypes(target types.Type) error {
	switch target.(type) {
	case *types.RefType:
	case *types.Void:
		node.TypeCheckPass.IsVoid = true
	default:
		return node.Errorf("cannot assign ref to %s", target)
	}

	node.TypeCheckPass.Type = target

	return nil
}

func (node *Ptr) CalculateStatistics(ctx context.Context) {
	node.Var.CalculateStatistics(ctx)
	node.Var.NameResolutionPass.Decl.StatisticsPass.PtrTaken = true
}

func (node *Ptr) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return nil, nil // this does not actually need the variable value
}

func (node *Ptr) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Ptr) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if !node.TypeCheckPass.IsVoid {
		p.AddI1RLit(bytecode.OPStackPtr, node.RegisterAllocationPass.Result.(bytecode.R), uint16(node.Var.NameResolutionPass.Decl.AddressResolutionPass.RelAddr), node.Position())
	}

	return nil
}
