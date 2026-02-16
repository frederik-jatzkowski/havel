package terminator

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block/instruction"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/types"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Conditional struct {
	tool.Node[Conditional]
	names.NameResolution[struct {
		IsMain     bool
		Then, Else function.Block
	}]

	Condition instruction.MemoryRead `parser:"'if':Keyword @@"`
	Then      string                 `parser:"'then':Keyword @Ident"`
	Else      string                 `parser:"'else':Keyword @Ident"`
}

var _ block.Terminator = (*Conditional)(nil)

func (node *Conditional) ResolveNames(ctx context.Context) error {
	thenTarget, err := contexttool.FromCtx[function.Block](ctx, node.Then)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Then = thenTarget

	elseTarget, err := contexttool.FromCtx[function.Block](ctx, node.Else)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.Else = elseTarget

	fn, err := contexttool.CurrentFromContext[*function.Function](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.IsMain = fn.Identifier() == names.SpecialMain

	return node.Condition.ResolveNames(ctx)
}

func (node *Conditional) ResolveTypes() error {
	actualType := node.Condition.Type()

	if !actualType.Equals(&types.Scalar{Size: 1}) {
		return node.Errorf("condition must be 1 byte but was %s", actualType)
	}

	return nil
}

func (node *Conditional) CalculateStatistics(ctx context.Context) {
	node.Condition.CalculateStatistics(ctx)
}

func (node *Conditional) AllocateRegisters(scope registeralloc.Scope) error {
	_, err := node.Condition.AllocateRegisters(scope)
	return err
}

func (node *Conditional) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if err := node.Condition.GenerateVirtualMachineAssembly(p); err != nil {
		return node.Wrap(err)
	}

	p.AddJumpToLabelIf(node.Condition.Register().(bytecode.R), node.NameResolutionPass.Then.FullyQualifiedIdentifier(), node.Position())
	p.AddJumpToLabel(node.NameResolutionPass.Else.FullyQualifiedIdentifier(), node.Position())

	return nil
}

func (node *Conditional) Successors() []controlflow.Node {
	return []controlflow.Node{
		node.NameResolutionPass.Then,
		node.NameResolutionPass.Else,
	}
}
