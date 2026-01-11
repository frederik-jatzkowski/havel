package terminator

import (
	"context"

	"github.com/frederik-jatzkowski/havel/pkg/hvil/architecture"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/program/function/block"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/runtime"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/lang/tool"
	"github.com/frederik-jatzkowski/havel/pkg/hvil/pass/registeralloc"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/assembly"
	"github.com/frederik-jatzkowski/havel/pkg/virtualmachine/bytecode"
)

type Return struct {
	tool.Node[Return]
	registeralloc.RegisterAllocation[struct {
		ExitCode architecture.Register
	}]

	Token string `parser:"@'return':Keyword" json:"-"`
}

var _ block.Terminator = (*Return)(nil)

func (node *Return) ResolveNames(_ context.Context) error {
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

func (node *Return) GenerateVirtualMachineAssembly(p *assembly.P, isMain bool) error {
	if isMain {
		// exit code 0
		p.AddLit(node.RegisterAllocationPass.ExitCode.(bytecode.R), 1, 0, node.Position())
		p.AddI1R(bytecode.OPExit, node.RegisterAllocationPass.ExitCode.(bytecode.R), node.Position())
	} else {
		//TODO implement me
		panic("implement me")
	}

	return nil
}

func (node *Return) Execute(vm *runtime.VirtualMachine) (*block.Block, error) {
	// there is no next block after a return statement
	return nil, nil
}
