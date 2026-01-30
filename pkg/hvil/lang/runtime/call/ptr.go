package call

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/types"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Ptr struct {
	tool.Node[Ptr]
	names.NameResolution[struct {
		Target *function.Function
	}]
	registeralloc.RegisterAllocation[struct {
		Result architecture.Register
	}]

	Name string `parser:"'ptr' '(' @Ident ')'"`
}

func (node *Ptr) ResolveNames(ctx context.Context) error {
	decl, err := contexttool.FromCtx[*function.Function](ctx, node.Name)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Target = decl

	return nil
}

func (node *Ptr) ResolveTypes(target types.Type) error {
	if err := target.CanBeAssignedDetailed(node.NameResolutionPass.Target.Signature()); err != nil {
		return node.Wrap(err)
	}

	return nil
}

func (node *Ptr) CalculateStatistics(_ context.Context) {
	node.NameResolutionPass.Target.StatisticsPass.AddressTaken++
}

func (node *Ptr) AllocateRegisters(_ registeralloc.Scope) ([]architecture.Register, error) {
	return nil, nil
}

func (node *Ptr) SetResultRegister(r architecture.Register) {
	node.RegisterAllocationPass.Result = r
}

func (node *Ptr) GenerateVirtualMachineAssembly(p *assembly.P) error {
	fqn := node.NameResolutionPass.Target.NameResolutionPass.Entry.FullyQualifiedIdentifier()

	p.AddLoadLabel(node.RegisterAllocationPass.Result.(bytecode.R), fqn, node.Position())

	return nil
}
