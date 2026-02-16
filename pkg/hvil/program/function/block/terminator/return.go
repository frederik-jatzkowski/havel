package terminator

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/architecture/virtualmachine/bytecode"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/optimization/controlflow"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/internal/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/tool"
	"github.com/frederik-jatzkowski/havel/pkg/tool/contexttool"
)

type Return struct {
	tool.Node[Return]
	names.NameResolution[struct {
		IsMain   bool
		Function *function.Function
	}]
	registeralloc.RegisterAllocation[struct {
		Temp architecture.Register
	}]

	Token string `parser:"@'return':Keyword" json:"-"`
}

var _ block.Terminator = (*Return)(nil)

func (node *Return) ResolveNames(ctx context.Context) error {
	fn, err := contexttool.CurrentFromContext[*function.Function](ctx)
	if err != nil {
		return node.Wrap(err)
	}

	node.NameResolutionPass.IsMain = fn.Identifier() == names.SpecialMain
	node.NameResolutionPass.Function = fn

	return nil
}

func (node *Return) ResolveTypes() error {
	return nil
}

func (node *Return) CalculateStatistics(ctx context.Context) {

}

func (node *Return) AllocateRegisters(scope registeralloc.Scope) error {
	r, ok := scope.GetScratchRegister()
	if !ok {
		return node.Errorf("failed to obtain exit code register")
	}

	scope.ReturnScratchRegisters(r)
	node.RegisterAllocationPass.Temp = r

	return nil
}

func (node *Return) GenerateVirtualMachineAssembly(p *assembly.P) error {
	temp := node.RegisterAllocationPass.Temp.(bytecode.R)
	if node.NameResolutionPass.IsMain {
		// exit code 0
		p.AddLit(temp, 1, 0, node.Position())
		p.AddI1R(bytecode.OPExit, temp, node.Position())
	} else {
		p.AddI1RLit(bytecode.OPStackPtr, temp, 0, node.Position())
		p.AddI2R(bytecode.OPLoad64, bytecode.PC, temp, node.Position())
	}

	return nil
}

func (node *Return) Successors() []controlflow.Node {
	return nil
}
