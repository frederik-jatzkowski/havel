package terminator

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool/contexttool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/names"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Return struct {
	tool.Node[Return]
	names.NameResolution[struct {
		IsMain bool
	}]
	registeralloc.RegisterAllocation[struct {
		ExitCode architecture.Register
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

	return nil
}

func (node *Return) ResolveTypes() error {
	return nil
}

func (node *Return) AllocateRegisters(arch architecture.Architecture) error {
	r, ok := arch.GetScratchRegister()
	if !ok {
		return node.Errorf("failed to obtain exit code register")
	}

	arch.ReturnScratchRegisters(r)
	node.RegisterAllocationPass.ExitCode = r

	return nil
}

func (node *Return) GenerateVirtualMachineAssembly(p *assembly.P) error {
	if node.NameResolutionPass.IsMain {
		// exit code 0
		p.AddLit(node.RegisterAllocationPass.ExitCode.(bytecode.R), 1, 0, node.Position())
		p.AddI1R(bytecode.OPExit, node.RegisterAllocationPass.ExitCode.(bytecode.R), node.Position())
	} else {
		p.AddI1RLit(bytecode.OPLoadStack64, bytecode.PC, 0, node.Position())
	}

	return nil
}

func (node *Return) Execute(vm *runtime.VirtualMachine) (function.Block, error) {
	// there is no next block after a return statement
	return nil, nil
}
