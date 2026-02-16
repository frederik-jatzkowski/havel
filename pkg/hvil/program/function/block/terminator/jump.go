package terminator

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Jump struct {
	tool.Node[Jump]
	names.NameResolution[struct {
		IsMain   bool
		Target   function.Block
		Function *function.Function
	}]

	Target string `parser:"'goto':Keyword @Ident"`
}

var _ block.Terminator = (*Jump)(nil)

func (node *Jump) ResolveNames(ctx context.Context) error {
	target, err := contexttool.FromCtx[function.Block](ctx, node.Target)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Target = target

	fn, err := contexttool.CurrentFromContext[*function.Function](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Function = fn
	node.NameResolutionPass.IsMain = fn.Identifier() == names.SpecialMain

	return nil
}

func (node *Jump) ResolveTypes() error {
	return nil
}

func (node *Jump) CalculateStatistics(ctx context.Context) {}

func (node *Jump) AllocateRegisters(scope registeralloc.Scope) error {
	return nil
}

func (node *Jump) GenerateVirtualMachineAssembly(p *assembly.P) error {
	p.AddJumpToLabel(node.NameResolutionPass.Target.FullyQualifiedIdentifier(), node.Position())

	return nil
}

func (node *Jump) Successors() []controlflow.Node {
	return []controlflow.Node{node.NameResolutionPass.Target}
}
