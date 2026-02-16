package call

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
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
	if err := target.EqualsDetailed(node.NameResolutionPass.Target.Signature()); err != nil {
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
