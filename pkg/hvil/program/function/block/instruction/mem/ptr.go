package mem

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/typecheck"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
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

	Var *instruction.VarRead `parser:"'ptr' '(' @@ ')'"`
}

func (node *Ptr) ResolveNames(ctx context.Context) error {
	return node.Var.ResolveNames(ctx)
}

func (node *Ptr) ResolveTypes(target types.Type) error {
	switch target.(type) {
	case *types.Ref:
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
	node.Var.NameResolutionPass.Decl.SetPtrTaken()
}

func (node *Ptr) AllocateRegisters(scope registeralloc.Scope) ([]architecture.Register, error) {
	return nil, nil // this does not actually need the variable value
}

func (node *Ptr) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Ptr) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if !node.TypeCheckPass.IsVoid {
		node.Var.NameResolutionPass.Decl.AddBytecodeVirtualmachinePtrInstruction(p, node.RegisterAllocationPass.Result.(bytecode.R))
	}

	return nil
}
